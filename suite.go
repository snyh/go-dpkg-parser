package dpkg

import (
	"path"
)

type Suite struct {
	Packages map[string]map[string]ControlFile

	Architecutres []string

	CodeName   string
	dataDir    string
	host       string
	limitArchs []string
}

func NewSuite(url string, codename string, dataDir string, archs ...string) (*Suite, error) {
	s := &Suite{
		limitArchs: archs,
		host:       url,
		CodeName:   codename,
		dataDir:    dataDir,
	}
	return s, s.build()
}

func (s Suite) FindBinaryBySource(sp SourcePackage, arch string) []BinaryPackage {
	var ret []BinaryPackage
	for _, name := range sp.GetBinary(arch) {
		b, err := s.FindBinary(name, arch)
		if err != nil {
			DebugPrintf("W: FindBinaryBySource(%s,%s)->%q: %v\n", sp.Package, arch, name, err)
			continue
		}
		ret = append(ret, b)
	}
	return ret
}

func (s Suite) FindBinary(name string, arch string) (BinaryPackage, error) {
	if arch == "all" {
		arch = s.limitArchs[0]
		if arch == "" {
			arch = s.Architecutres[0]
		}
	}
	r, ok := s.FindControl(name, arch)
	if !ok {
		return BinaryPackage{}, NotFoundError{"Architecutre of " + string(arch)}
	}
	return r.ToBinary()
}

func (s Suite) FindControl(name string, arch string) (ControlFile, bool) {
	db, ok := s.Packages[arch]
	if !ok {
		return ControlFile{}, false
	}
	r, ok := db[name]
	return r, ok
}

func (s Suite) FindSource(name string) (SourcePackage, error) {
	r, ok := s.FindControl(name, "source")
	if !ok {
		return SourcePackage{}, NotFoundError{name}
	}
	return r.ToSource()
}

func (s *Suite) tryDownload() (ReleaseFile, error) {
	rfPath := path.Join(s.dataDir, ReleaseFileName)

	rf, err := DownloadReleaseFile(s.host, s.CodeName, rfPath)
	if err != nil {
		return rf, err
	}
	if len(s.limitArchs) != 0 {
		rf.Architectures = (UnionSet(rf.Architectures, s.limitArchs))
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

	fs := make(map[string][]string)
	for _, f := range rf.FileInfos() {
		fs[f.Architecture] = append(fs[f.Architecture], path.Join(s.dataDir, f.Path))
	}

	s.Packages = make(map[string]map[string]ControlFile)
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
