package dpkg

import (
	"fmt"
	"path"
	"strings"
	"sync"
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
	rf, err := DownloadReleaseFile(s.host, s.CodeName, path.Join(s.dataDir, ReleaseFileName))
	if err != nil {
		return err
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

type Suite struct {
	cacheDir string
	name     string
	index    *PackageDBIndex
	dbs      map[Architecture]PackageDB
	repoURL  string

	dbLock sync.Mutex
}

func NewSuite(cacheDir string, repoURL string, name string) (*Suite, error) {
	if repoURL == "" || cacheDir == "" || name == "" {
		return nil, fmt.Errorf("Please setup packages.newSuite")
	}
	s := &Suite{
		cacheDir: cacheDir,
		name:     name,
		dbs:      make(map[Architecture]PackageDB),
		repoURL:  repoURL,
	}
	return s, nil
}

func (m *Suite) Search(q string) []string {
	var r = make(map[string]struct{})
	for _, data := range m.index.PackagePaths {
		for id := range data {
			if strings.Contains(id, q) {
				r[id] = struct{}{}
			}
		}
	}
	return sortMapString(r)
}

func (m *Suite) QueryPath(id string, arch Architecture) (string, bool) {
	data, ok := m.index.PackagePaths[arch]
	if !ok {
		return "", false
	}
	path, ok := data[id]
	return path, ok
}

func (m *Suite) Get(id string) (BinaryPackage, bool) {
	archs := m.index.PackageArchitectures(id)
	for _, arch := range archs {
		DB, _ := m.getDB(arch)
		t, ok := DB[id]
		if !ok {
			continue
		}
		t.Architectures = archs
		return t, true
	}
	return BinaryPackage{}, false
}

func (m *Suite) getDB(arch Architecture) (PackageDB, error) {
	// If we don't lock this, the loadPackageDB maybe invoked too many times
	// to cause memory exploded.
	m.dbLock.Lock()
	defer m.dbLock.Unlock()

	DB, ok := m.dbs[arch]

	if !ok {
		var err error
		DB, err = loadPackageDB(buildDBPath(m.cacheDir, m.name, DBName(arch)))
		if err != nil {
			return nil, err
		}
		m.dbs[arch] = DB
	}
	return DB, nil
}
