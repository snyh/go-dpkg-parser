package dpkg

type Archive struct {
	Virtuals     map[string][]string
	Packages     map[string]ControlFile
	Architecture string

	cache map[string]error
}

func NewArchive(arch string) Archive {
	return Archive{
		Architecture: arch,
		Packages:     make(map[string]ControlFile),
		Virtuals:     make(map[string][]string),
		cache:        make(map[string]error),
	}
}

func (a Archive) FindControl(name string) (ControlFile, bool) {
	cf, ok := a.Packages[name]
	return cf, ok
}

func (a Archive) FindProvider(name string) []string {
	return a.Virtuals[name]
}

func (a Archive) FindSource(name string) (SourcePackage, error) {
	r, ok := a.Packages[name]
	if !ok {
		return SourcePackage{}, NotFoundError{name}
	}
	return r.ToSource()
}

func (a Archive) FindBinary(name string) (BinaryPackage, error) {
	r, ok := a.Packages[name]
	if !ok {
		return BinaryPackage{}, NotFoundError{name}
	}
	return r.ToBinary()
}

func (a Archive) SatisfyDepends(str string) error {
	panic("Not Implement")
}
