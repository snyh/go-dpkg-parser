package dpkg

import (
	"fmt"
)

type Suite struct {
	Packages map[string]map[string]ControlFile

	Virtuals map[string]map[string][]string

	Architecutres []string

	CodeName string

	dataDir    string
	host       string
	limitArchs []string
	hash       string
}

func NewSuite(url string, codename string, dataDir string, hash string, archs ...string) (*Suite, error) {
	s := &Suite{
		limitArchs: archs,
		host:       url,
		CodeName:   codename,
		dataDir:    dataDir,
		hash:       hash,
	}
	return s, s.build()
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

func (s Suite) FindBinaryAny(name string, archs []string) (BinaryPackage, error) {
	for _, arch := range archs {
		if refs, ok := s.Virtuals[arch][name]; ok && len(refs) != 0 {
			// TODO: Design properly API for handling multiple virtual packages
			name = refs[0]
		}
		r, ok := s.FindControl(name, arch)
		if ok {
			return r.ToBinary()
		}
	}
	return BinaryPackage{}, NotFoundError{fmt.Sprintf("Architecutre of %v(%d) for %q", archs, len(archs), name)}
}

func (s Suite) FindBinary(name string, arch string) (BinaryPackage, error) {
	var archs = []string{arch}
	if arch == "all" {
		if len(s.limitArchs) != 0 {
			archs = UnionSet(s.limitArchs, s.Architecutres)
		} else {
			archs = s.Architecutres
		}
	}
	return s.FindBinaryAny(name, archs)
}
