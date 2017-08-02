package dpkg

import (
	"bytes"
	"fmt"
	"path"
)

func LoadPackages(fPath string) ([]ControlFile, error) {
	f, err := ReadFile(fPath, IsGzip(fPath))
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return NewControlFiles(f, ScanBufferSize)
}

func DownloadReleaseFile(repoURL string, codeName string) (ControlFile, error) {
	var r ControlFile
	url := fmt.Sprintf("%s/dists/%s/%s", repoURL, codeName, ReleaseFileName)
	buf := bytes.NewBuffer(nil)
	err := DownloadTo(url, buf)
	if err != nil {
		return r, err
	}
	return NewControlFile(string(buf.Bytes()))
}

func (s *Suite) downloadReleaseFile() (ReleaseFile, error) {
	var rf ReleaseFile
	raw, err := DownloadReleaseFile(s.host, s.CodeName)
	if err != nil {
		return rf, err
	}
	rf, err = raw.ToReleaseFile()
	if err != nil {
		return rf, err
	}

	// Update Suite Hash
	s.hash = rf.Hash

	rPath := s.rootDir(ReleaseFileName)

	err = WriteToFile([]byte(raw.Raw), rPath, 0644)
	if err != nil {
		return rf, err
	}
	return rf, err
}

func (s *Suite) loadReleaseFile() (ReleaseFile, error) {
	if s.hash == "" {
		return s.downloadReleaseFile()
	}
	rPath := s.rootDir(ReleaseFileName)
	rf, err := LoadReleaseFile(rPath)
	if rf.Hash != s.hash {
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
		rf.Architectures = (UnionSet(rf.Architectures, s.limitArchs))
	}
	return rf, DownloadRepository(s.host, rf, s.rootDir())
}

func (s *Suite) rootDir(subPath ...string) string {
	if s.hash == "" {
		panic("Empty　Hash")
	}
	root := path.Join(s.dataDir, s.hash)
	return path.Join(append([]string{root}, subPath...)...)
}

func (s *Suite) build() error {
	rf, err := s.prepareDownload()
	if err != nil {
		return err
	}

	s.Architecutres = rf.Architectures

	fs := make(map[string][]string)
	hashs := make(map[string][]string)
	for _, f := range rf.FileInfos() {
		fs[f.Architecture] = append(fs[f.Architecture], s.rootDir(f.Path))
		hashs[f.Architecture] = append(hashs[f.Architecture], f.MD5)
	}

	s.Packages = make(map[string]map[string]ControlFile)
	s.Virtuals = make(map[string]map[string][]string)
	for arch, srcs := range fs {
		cache, err := loadOrBuildCache(
			path.Join(s.dataDir, "db", HashArrayString(hashs[arch])),
			srcs...)
		if err != nil {
			return err
		}
		s.Packages[arch], s.Virtuals[arch] = cache.Pkgs, cache.Virtual
	}
	return err
}

type pkgCache struct {
	Pkgs    map[string]ControlFile
	Virtual map[string][]string
}

func loadOrBuildCache(cacheFile string, files ...string) (pkgCache, error) {
	var cache = pkgCache{
		Pkgs:    make(map[string]ControlFile),
		Virtual: make(map[string][]string),
	}
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
			cache.Pkgs[name] = cf
			for _, p := range cf.GetArray("provides", ",") {
				cache.Virtual[p] = append(cache.Virtual[p], name)
			}
		}
	}
	return cache, storeGOB(cacheFile, cache)
}

// DownloadRepository download files from rf.FileInfos()
func DownloadRepository(repoURL string, rf ReleaseFile, rootDir string) error {
	for _, f := range rf.FileInfos() {
		url := repoURL + "/dists/" + rf.CodeName + "/" + f.Path
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
