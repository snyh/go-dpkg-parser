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

func (s Suite) FindProvider(name string, arch string) []string {
	return s.Virtuals[arch][name]
}

func (s Suite) FindBinary(name string, arch string) (BinaryPackage, error) {
	if refs := s.FindProvider(name, arch); len(refs) != 0 {
		// TODO: Design properly API for handling multiple virtual packages
		name = refs[0]
	}
	r, ok := s.FindControl(name, arch)
	if ok {
		return r.ToBinary()
	}
	return BinaryPackage{}, NotFoundError{fmt.Sprintf("Binary %q with architecutre %q", name, arch)}
}
