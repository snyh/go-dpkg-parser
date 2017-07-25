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
	SourcePackage map[string]ControlFile
	Sections      []string
	CodeName      string
	Architecutres []Architecture

	dataDir string
	host    string
}

func NewSuite1(url string, codename string, dataDir string) (*Suite1, error) {
	s := &Suite1{
		host:     url,
		CodeName: codename,
		dataDir:  dataDir,
	}
	return s, s.build()
}

func (s *Suite1) build() error {
	rf, err := DownloadReleaseFile(s.host, s.CodeName, path.Join(s.dataDir, ReleaseFileName))
	if err != nil {
		return err
	}
	_, err = DownloadRepository(s.host, rf, s.dataDir)

	// TODO: Build the Caches..
	return err
}

func (s *Suite1) ListSource() {
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
