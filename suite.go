package dpkg

import (
	"fmt"
	"path"
)

func DownloadReleaseFile(repoURL string, codeName string, fpath string) (ReleaseFile, error) {
	var r ReleaseFile
	url := fmt.Sprintf("%s/dists/%s/%s", repoURL, codeName, ReleaseFileName)

	// download Release File
	f, err := DownloadAndOpen(url, fpath)
	if err != nil {
		return r, fmt.Errorf("DownloadReleaseFile  http.Get(%q) failed:(%v)", url, err)
	}
	defer f.Close()

	// build Release File
	cf, err := NewControlFile(f)
	if err != nil {
		return r, fmt.Errorf("DownloadReleaseFile invalid Release file(%q) : %v", url, err)
	}
	return cf.ToReleaseFile()
}

type Suite1 struct {
	BinaryPackage map[Architecture]map[string]ControlFile

	Architecutres []Architecture

	CodeName string
	dataDir  string
	host     string
}

func NewSuite1(url string, codename string, dataDir string) (*Suite1, error) {
	s := &Suite1{
		host:     url,
		CodeName: codename,
		dataDir:  dataDir,
	}
	return s, s.build()
}

func (s *Suite1) FindBinary(name string) (BinaryPackage, error) {
	for arch, db := range s.BinaryPackage {
		if arch == "source" {
			continue
		}
		return db[name].ToBinary()
	}
	return BinaryPackage{}, nil
}

func (s *Suite1) build() error {
	rfPath := path.Join(s.dataDir, ReleaseFileName)
	oldRF, _ := GetReleaseFile(rfPath)

	rf, err := DownloadReleaseFile(s.host, s.CodeName, rfPath)
	if err != nil {
		return err
	}

	if oldRF.Hash() == rf.Hash() {
		return nil
	}

	_, err = DownloadRepository(s.host, rf, s.dataDir)

	s.Architecutres = rf.Architectures

	fs := make(map[Architecture][]string)
	for _, f := range rf.FileInfos() {
		fs[f.Architecture] = append(fs[f.Architecture], path.Join(s.dataDir, f.Path))
	}

	s.BinaryPackage = make(map[Architecture]map[string]ControlFile)
	for arch, cs := range fs {
		s.BinaryPackage[arch], err = buildCache(path.Join(s.dataDir, "cache-"+string(arch)+".db"), cs...)
		if err != nil {
			return err
		}
	}
	return err
}

func buildCache(cacheFile string, files ...string) (map[string]ControlFile, error) {
	r := make(map[string]ControlFile)
	if err := loadGOB(cacheFile, &r); err == nil {
		return r, nil
	}

	for _, f := range files {
		cfs, err := LoadControlFileGroup(f)
		if err != nil {
			return nil, err
		}
		for _, cf := range cfs {
			r[cf.GetString("Package")] = cf
		}
	}
	return r, storeGOB(cacheFile, r)
}

func (s *Suite1) ListSource() []string {
	var re []string
	for p := range s.BinaryPackage["source"] {
		re = append(re, p)
	}
	return re
}
