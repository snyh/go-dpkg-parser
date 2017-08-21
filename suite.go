package dpkg

import (
	"bytes"
	"fmt"
	"path"
)

type Suite struct {
	Archives map[string]Archive

	Suite string
	Host  string
	Hash  string

	dataDir    string
	limitArchs []string
}

func NewSuite(url string, suite string, dataDir string, hash string, archs ...string) (*Suite, error) {
	s := &Suite{
		Archives:   make(map[string]Archive),
		limitArchs: archs,
		Host:       url,
		Suite:      suite,
		dataDir:    dataDir,
		Hash:       hash,
	}
	return s, s.build()
}

func LoadPackages(fPath string) ([]ControlFile, error) {
	f, err := ReadFile(fPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewControlFiles(f, ScanBufferSize)
}

func DownloadReleaseFile(repoURL string, suiteName string) (ControlFile, error) {
	var r ControlFile
	url := fmt.Sprintf("%s/dists/%s/%s", repoURL, suiteName, ReleaseFileName)
	buf := bytes.NewBuffer(nil)
	err := DownloadTo(url, buf)
	if err != nil {
		return r, err
	}
	return NewControlFile(string(buf.Bytes()))
}

func (s *Suite) downloadReleaseFile() (ReleaseFile, error) {
	var rf ReleaseFile
	raw, err := DownloadReleaseFile(s.Host, s.Suite)
	if err != nil {
		return rf, err
	}
	rf, err = raw.ToReleaseFile()
	if err != nil {
		return rf, err
	}

	// Update Suite Hash
	s.Hash = rf.Hash

	rPath := s.rootDir(ReleaseFileName)

	err = WriteToFile([]byte(raw.Raw), rPath, 0644)
	if err != nil {
		return rf, err
	}
	return rf, err
}

func (s *Suite) loadReleaseFile() (ReleaseFile, error) {
	if s.Hash == "" {
		return s.downloadReleaseFile()
	}
	rPath := s.rootDir(ReleaseFileName)
	rf, err := LoadReleaseFile(rPath)
	if rf.Hash != s.Hash {
		DebugPrintf("Invalid cache release file %q. Redownload it\n", rPath)
		return s.downloadReleaseFile()
	}
	DebugPrintf("Loaded releae file %q from cache", rPath)
	return rf, err
}

func (s *Suite) prepareDownload() (ReleaseFile, error) {
	rf, err := s.loadReleaseFile()
	if err != nil {
		return rf, err
	}
	if len(s.limitArchs) != 0 {
		rf.Architectures = (IntersectionSet(rf.Architectures, s.limitArchs))
	}

	rf.Suite = s.Suite // workaround insane suite value
	return rf, DownloadRepository(s.Host, rf, s.rootDir())
}

func (s *Suite) rootDir(subPath ...string) string {
	if s.Hash == "" {
		panic("Empty　Hash")
	}
	root := path.Join(s.dataDir, s.Hash)
	return path.Join(append([]string{root}, subPath...)...)
}

func (s *Suite) build() error {
	rf, err := s.prepareDownload()
	if err != nil {
		return err
	}

	contents := make(map[string][]string)
	fs := make(map[string][]string)
	hashs := make(map[string][]string)
	for _, f := range rf.FileInfos() {
		switch f.Type {
		case tCONTROLFILES:
			fs[f.Architecture] = append(fs[f.Architecture], s.rootDir(f.Path))
		case tCONTENTS:
			contents[f.Architecture] = append(contents[f.Architecture], s.rootDir(f.Path))
		default:
			DebugPrintf("Unknown component type %q for %v\n", f.Type, f.Path)
		}
		hashs[f.Architecture] = append(hashs[f.Architecture], f.MD5)
	}

	for arch, srcs := range fs {
		cache, err := loadOrBuildArchive(
			arch,
			path.Join(s.dataDir, "db", HashArrayString(hashs[arch])),
			srcs...)
		if err != nil {
			return err
		}
		s.Archives[arch] = cache
	}
	return err
}

func loadOrBuildArchive(arch string, cacheFile string, files ...string) (Archive, error) {
	cache := NewArchive(arch)
	if err := EnsureDirectory(path.Dir(cacheFile)); err != nil {
		return cache, err
	}
	err := loadGOB(cacheFile, &cache)
	if err == nil {
		return cache, nil
	}

	DebugPrintf("Build %q from %v\n", cacheFile, files)

	for _, fpath := range files {
		cfs, err := LoadPackages(fpath)
		if err != nil {
			return cache, err
		}
		for _, cf := range cfs {
			name := cf.Get("Package")
			cache.Packages[name] = cf
			for _, p := range parseProvides(cf.Get("provides")) {
				cache.Virtuals[p] = append(cache.Virtuals[p], name)
			}
		}
	}
	return cache, storeGOB(cacheFile, cache)
}

// DownloadRepository download files from rf.FileInfos()
func DownloadRepository(repoURL string, rf ReleaseFile, rootDir string) error {
	for _, f := range rf.FileInfos() {
		url := repoURL + "/dists/" + rf.Suite + "/" + f.Path
		target := path.Join(rootDir, f.Path)
		hash, _ := HashFile(target)
		if hash == f.MD5 {
			DebugPrintf("%q to %q is cached\n", url, target)
			continue
		}

		err := DownloadToFile(url, target)
		if err != nil {
			return err
		}
	}
	return nil
}
