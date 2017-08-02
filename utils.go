package dpkg

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sort"
	"strings"
	"unicode"
)

func WriteToFile(content []byte, target string, mode os.FileMode) error {
	err := EnsureDirectory(path.Dir(target))
	if err != nil {
		return err
	}
	return ioutil.WriteFile(target, content, mode)
}

func IsGzip(fpath string) bool {
	return strings.HasSuffix(strings.ToLower(fpath), ".gz")
}

type readCloserWrap struct {
	*gzip.Reader
	a io.ReadWriteCloser
}

func (w readCloserWrap) Close() error {
	w.a.Close()
	return w.Reader.Close()
}

func ReadFile(fpath string, unzip bool) (io.ReadCloser, error) {
	f, err := os.Open(fpath)
	if err != nil {
		return nil, err
	}
	if !unzip {
		return f, err
	}
	gr, err := gzip.NewReader(f)
	if err != nil {
		return nil, FormatError{"LoadControlFileGroup", fpath, err}
	}
	return readCloserWrap{gr, f}, nil
}

func DownloadTo(url string, w io.Writer) error {
	reps, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("can't download %q : %v", url, err)
	}
	defer reps.Body.Close()
	if reps.StatusCode != 200 {
		return fmt.Errorf("can't download %q : %v", url, reps.Status)
	}
	_, err = io.Copy(w, reps.Body)
	return err
}

// download download the url content to "dest" file.
func DownloadToFile(url string, dest string) error {
	DebugPrintf("Downloading %q to %q\n", url, dest)
	if err := EnsureDirectory(path.Dir(dest)); err != nil {
		return err
	}
	f, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("Can't create file %s", url)
	}
	defer f.Close()
	return DownloadTo(url, f)
}

func UnionSet(s1, s2 []string) []string {
	var ret []string
	for _, i := range s1 {
		for _, j := range s2 {
			if i == j {
				ret = append(ret, i)
				break
			}
		}
	}
	return ret
}

func HashFiles(fpaths ...string) (string, error) {
	var hashs []string
	for _, f := range fpaths {
		v, err := HashFile(f)
		if err != nil {
			return "", err
		}
		hashs = append(hashs, v)
	}
	return HashArrayString(hashs), nil
}
func HashArrayString(s []string) string {
	sort.Strings(s)
	var r string
	for _, h := range s {
		r += h
	}
	return HashBytes([]byte(r))
}

func HashBytes(bs []byte) string {
	hash := md5.New()
	hash.Write(bs)
	var r [16]byte
	copy(r[:], hash.Sum(nil))
	return fmt.Sprintf("%x", r)
}

func HashFile(fpath string) (string, error) {
	bs, err := ioutil.ReadFile(fpath)
	if err != nil {
		return "", err
	}
	return HashBytes(bs), nil
}

func TrimLeftSpace(d []byte) []byte {
	return bytes.TrimFunc(d, unicode.IsSpace)
}

func EnsureDirectory(t string) error {
	s, err := os.Stat(t)
	if err != nil {
		err := os.MkdirAll(t, 0755)
		if err != nil {
			DebugPrintf("EnsureDirectory failed: %v", err)
		}
		return err
	} else {
		if !s.IsDir() {
			return fmt.Errorf("%q is a regular file", t)
		}
	}
	return nil
}

func sortMapString(d map[string]struct{}) []string {
	var r = make([]string, 0)
	for k := range d {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}

func loadGOB(fpath string, obj interface{}) error {
	f, err := os.Open(fpath)
	if err != nil {
		return fmt.Errorf("store %q failed --> %v.", fpath, err)
	}
	defer f.Close()
	return gob.NewDecoder(f).Decode(obj)
}
func storeGOB(fpath string, obj interface{}) error {
	f, err := os.Create(fpath)
	if err != nil {
		return fmt.Errorf("store %q failed --> %v.", fpath, err)
	}
	defer f.Close()
	return gob.NewEncoder(f).Encode(obj)
}
