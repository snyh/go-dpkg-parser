package dpkg

import (
	"bytes"
	"crypto/md5"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"unicode"
)

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
