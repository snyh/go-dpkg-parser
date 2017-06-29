package dpkg

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
)

func loadPackagesaDBIndex(fpath string) (*PackageDBIndex, error) {
	index := &PackageDBIndex{}
	err := loadGOB(fpath, &index)
	return index, err
}

func loadPackageDB(fpath string) (PackageDB, error) {
	var obj = make(PackageDB)
	err := loadGOB(fpath, &obj)
	return obj, err
}

func (m *Suite) DataHash() string {
	rf, err := GetReleaseFile(m.cacheDir, m.name)
	if err != nil {
		panic(err)
	}
	return rf.Hash()
}

func (m *Suite) loadIndex(targetDir string) error {
	index, err := loadPackagesaDBIndex(buildDBPath(targetDir, m.name, DBIndexName))
	if err != nil {
		return fmt.Errorf("UpdateDB: failed load db index: %v", err)
	}
	m.index = index
	m.dbs = make(map[Architecture]PackageDB)
	return nil
}

func (m *Suite) UpdateDB() error {
	// always trying load index file at the end, no matter the result.
	defer m.loadIndex(m.cacheDir)

	// 1. comparing ReleaseFile.Hash to check whether need update
	os.MkdirAll(m.cacheDir, 0755)
	targetDir, err := ioutil.TempDir(m.cacheDir, "lastore.packages.partition")
	if err != nil {
		return fmt.Errorf("UpdateDB: failed create temp directory: %v", err)
	}
	defer os.RemoveAll(targetDir)

	rf, err := DownloadReleaseFile(m.repoURL, m.name, targetDir)
	if err != nil {
		return fmt.Errorf("UpdateDB: failed download release file: %v", err)
	}

	h := m.DataHash()
	if h == rf.Hash() {
		changed, err := DownloadRepository(m.repoURL, rf, m.cacheDir)
		if err != nil {
			return err
		}
		if changed {
			BuildCache(rf, m.cacheDir, m.cacheDir)
		}
		m.loadIndex(m.cacheDir)
		return nil
	}

	// 2. download data and build DBs to tmp/${xx} directory
	_, err = DownloadRepository(m.repoURL, rf, targetDir)
	if err != nil {
		return fmt.Errorf("UpdateDB: failed download repository files: %v", err)
	}
	err = BuildCache(rf, targetDir, targetDir)
	if err != nil {
		return fmt.Errorf("UpdateDB: failed build dbs: %v", err)
	}

	// 3. unload index
	m.index = nil

	// 4. removing old data on filesystem.
	os.RemoveAll(buildDBPath(m.cacheDir, m.name))

	// 5. moving tmp/${xx} to $cacheDir on filesystem.
	return os.Rename(buildDBPath(targetDir, m.name), buildDBPath(m.cacheDir, m.name))
}

// DownloadRepository download files from rf.FileInfos()
// it ignoring unchanged file by checking MD5 value.
// return whether changed and error if any.
func DownloadRepository(repoURL string, rf ReleaseFile, targetDir string) (bool, error) {
	changed := false
	for _, f := range rf.FileInfos() {
		url := repoURL + "/dists/" + rf.CodeName + "/" + f.Path
		target := path.Join(targetDir, rf.CodeName, "raw", f.Path)
		if HashFile(target) == f.MD5 {
			continue
		}
		changed = true

		err := download(url, target, f.Gzip)
		if err != nil {
			return changed, err
		}
	}
	return changed, nil
}

func DownloadReleaseFile(repoURL string, codeName string, targetDir string) (ReleaseFile, error) {
	url := fmt.Sprintf("%s/dists/%s/Release", repoURL, codeName)
	reps, err := http.Get(url)
	if err != nil {
		return ReleaseFile{}, fmt.Errorf("DownloadReleaseFile  http.Get(%q) failed:(%v)", url, err)
	}
	defer reps.Body.Close()
	// There is no need check reps.StatusCode, ReadReleaseFile will check the validation.

	// download Release File
	fpath := path.Join(targetDir, codeName, "Release")
	os.MkdirAll(path.Dir(fpath), 0755)
	f, err := os.Create(fpath)
	if err != nil {
		return ReleaseFile{}, fmt.Errorf("DownloadReleaseFile(%q) can't create file %q : %v", url, fpath, err)
	}
	defer f.Close()

	// build Release File
	bs, _ := ioutil.ReadAll(io.TeeReader(reps.Body, f))
	cf, err := NewControlFile(bs)
	if err != nil {
		return ReleaseFile{}, fmt.Errorf("DownloadReleaseFile invalid Release file(%q) : %v", url, err)
	}
	return cf.ToReleaseFile()
}

// DownloadTee download the url content to "dest" file
// and return a io.Reader to read it.
func download(url string, dest string, gz bool) error {
	os.MkdirAll(path.Dir(dest), 0755)

	reps, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("can't download %q : %v", url, err)
	}
	defer reps.Body.Close()

	if reps.StatusCode != 200 {
		return fmt.Errorf("can't download %q : %v", url, reps.Status)
	}

	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("Can't create file %s", url)
	}
	defer f.Close()

	n, err := io.Copy(f, reps.Body)
	if err != nil {
		return fmt.Errorf("DownloadTo: write content(%d) failed:%v", n, err)
	}

	return nil
}
