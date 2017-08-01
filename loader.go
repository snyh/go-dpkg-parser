package dpkg

import (
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func DownloadReleaseFile(repoURL string, codeName string) (ReleaseFile, error) {
	var r ReleaseFile
	url := fmt.Sprintf("%s/dists/%s/%s", repoURL, codeName, ReleaseFileName)

	reps, err := http.Get(url)
	if err != nil {
		return r, fmt.Errorf("can't download %q : %v", url, err)
	}
	defer reps.Body.Close()

	if reps.StatusCode != 200 {
		return r, fmt.Errorf("can't download %q : %v", url, reps.Status)
	}
	bs, err := ioutil.ReadAll(reps.Body)
	if err != nil {
		return r, fmt.Errorf("can't read content %q : %v", url, err)
	}

	// build Release File
	cf, err := NewControlFile(string(bs))
	if err != nil {
		return r, fmt.Errorf("DownloadReleaseFile invalid Release file(%q) : %v", url, err)
	}
	return cf.ToReleaseFile()
}

func LoadControlFileGroup(fPath string) ([]ControlFile, error) {
	f, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if strings.HasSuffix(strings.ToLower(fPath), ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return nil, FormatError{"LoadControlFileGroup", fPath, err}
		}
		defer gr.Close()
		return NewControlFiles(gr, ScanBufferSize)
	}
	return NewControlFiles(f, ScanBufferSize)
}

// DownloadRepository download files from rf.FileInfos()
// it ignoring unchanged file by checking MD5 value.
// return whether changed and error if any.
func DownloadRepository(repoURL string, rf ReleaseFile, rootDir string) (bool, error) {
	changed := false
	for _, f := range rf.FileInfos() {
		url := repoURL + "/dists/" + rf.CodeName + "/" + f.Path
		target := path.Join(rootDir, f.Path)
		hash, _ := HashFile(target)
		if hash == f.MD5 {
			DebugPrintf("%q to %q is cached\n", url, target)
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

// download download the url content to "dest" file.
// unpack only support gz file now.
func download(url string, dest string, unpack bool) error {
	DebugPrintf("Downloading %q to %q\n", url, dest)

	if err := EnsureDirectory(path.Dir(dest)); err != nil {
		return err
	}

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
