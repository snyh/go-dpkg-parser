package dpkg

import (
	"fmt"
	"io"
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
		return ""
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

// DownloadRepository download files from rf.FileInfos()
// it ignoring unchanged file by checking MD5 value.
// return whether changed and error if any.
func DownloadRepository(repoURL string, rf ReleaseFile, targetDir string) (bool, error) {
	changed := false
	for _, f := range rf.FileInfos() {
		url := repoURL + "/dists/" + rf.CodeName + "/" + f.Path
		target := path.Join(targetDir, f.Path)
		if HashFile(target) == f.MD5 {
			if Debug {
				fmt.Printf("%q to %q is cached\n", url, target)
			}
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

func BuildSuiteCache(rf ReleaseFile, cacheDir string) error {
	for _, f := range rf.FileInfos() {
		fmt.Printf("Caching %q to %q\n", f.Path, f.Path)
	}
	return nil
}

// download download the url content to "dest" file.
// unpack only support gz file now.
func download(url string, dest string, unpack bool) error {
	if Debug {
		fmt.Printf("Downloading %q to %q\n", url, dest)
	}
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
