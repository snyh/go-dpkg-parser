package dpkg

type Archive struct {
	Virtuals     map[string][]string
	Packages     map[string]ControlFile
	dependCache  map[string]DependStatus
	Architecture string
}

func (a Archive) DependCheck(name string) error {
	info, ok := a.dependCache[name]
	if ok {
		return info.err
	}
	return a.parseDepend(name)
}

func (a Archive) FindControl(name string) (ControlFile, bool) {
	cf, ok := a.Packages[name]
	return cf, ok
}

func (a Archive) FindProvider(name string, arch string) []string {
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
