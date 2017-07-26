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

type Suite struct {
	Packages map[Architecture]map[string]ControlFile

	Architecutres []Architecture

	CodeName string
	dataDir  string
	host     string
}

func NewSuite(url string, codename string, dataDir string) (*Suite, error) {
	s := &Suite{
		host:     url,
		CodeName: codename,
		dataDir:  dataDir,
	}
	return s, s.build()
}

func (s *Suite) FindBinary(name string, arch Architecture) (BinaryPackage, error) {
	if arch == "all" {
		arch = s.Architecutres[0]
	}
	for arch, db := range s.Packages {
		if arch == "source" {
			continue
		}
		return db[name].ToBinary()
	}
	return BinaryPackage{}, NotFoundError{"Architecutre of " + string(arch)}
}

func (s *Suite) tryDownload() (ReleaseFile, error) {
	rfPath := path.Join(s.dataDir, ReleaseFileName)
	oldRF, _ := GetReleaseFile(rfPath)

	rf, err := DownloadReleaseFile(s.host, s.CodeName, rfPath)
	if err != nil {
		return rf, err
	}

	if oldRF.Hash() == rf.Hash() {
		return rf, nil
	}

	_, err = DownloadRepository(s.host, rf, s.dataDir)
	return rf, err
}

func (s *Suite) build() error {
	rf, err := s.tryDownload()
	if err != nil {
		return err
	}

	s.Architecutres = rf.Architectures

	fs := make(map[Architecture][]string)
	for _, f := range rf.FileInfos() {
		fs[f.Architecture] = append(fs[f.Architecture], path.Join(s.dataDir, f.Path))
	}

	s.Packages = make(map[Architecture]map[string]ControlFile)
	for arch, cs := range fs {
		s.Packages[arch], err = buildCache(path.Join(s.dataDir, "cache-"+string(arch)+".db"), cs...)
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

func (s *Suite) FindBinaryBySource(sp SourcePackage, arch Architecture) []BinaryPackage {
	var ret []BinaryPackage
	for _, name := range sp.Binary {
		b, err := s.FindBinary(name, arch)
		if err != nil {
			continue
		}
		ret = append(ret, b)
	}
	return ret
}

func (s *Suite) FindSource(name string) (SourcePackage, error) {
	srcs, ok := s.Packages["source"]
	if !ok {
		return SourcePackage{}, NotFoundError{name}
	}
	r, ok := srcs[name]
	return r.ToSource()
}

func (s *Suite) ListSource() map[string]SourcePackage {
	var re = make(map[string]SourcePackage)
	for p, cf := range s.Packages["source"] {
		s, err := cf.ToSource()
		if err != nil {
			fmt.Println("W:", err)
			continue
		}
		re[p] = s
	}
	return re
}
