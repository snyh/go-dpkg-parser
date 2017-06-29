package dpkg

import (
	"encoding/gob"
	"fmt"
	"os"
	"sort"
)

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
