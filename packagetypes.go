package dpkg

import (
	"fmt"
	"path"
	"strconv"
	"strings"
)

type BinaryPackage struct {
	Package       string        `json:"package"`
	Source        string        `json:"source"`
	SourceVersion string        `json:"source_version"`
	Version       string        `json:"version"`
	InstalledSize int           `json:"installed_size"`
	Size          int           `json:"size"`
	Architectures Architectures `json:"architectures"`
	Description   string        `json:"description"`
	Filename      string        `json:"filename"`
	Tag           string        `json:"tag"`
	Homepage      string        `json:"homepage"`
	SHA256        string        `json:"sha256"`
	Maintainer    string        `json:"maintainer"`
	Priority      string        `json:"priority"`

	conflicts  string
	depends    string
	preDepends string
}

type PackageListItem struct {
	Name     string
	Ptype    string
	Section  string
	Priority string
	Archs    []string
	Profile  string
	Essional bool
}

func (item PackageListItem) Support(arch string) bool {
	if arch == "any" {
		panic("It's wrong to query package by architecture of any.")
	}
	for _, i := range item.Archs {
		switch i {
		case "any", "linux-any", "all":
			return true
		case arch, "any-" + arch:
			return true
		}
	}
	return false
}

type SourcePackage struct {
	Package      string            `json:"package"`
	Version      string            `json:"version"`
	Description  string            `json:"description"`
	Homepage     string            `json:"homepage"`
	Architecture []string          `json:"architecture"`
	Maintainer   string            `json:"maintainer"`
	Format       string            `json:"format"`
	Binary       []string          `json:"binary"`
	PackageList  []PackageListItem `json:"package_list"`

	Section  string `json:"section"`
	Priority string `json:"priority"`

	directory string
	files256  []string

	buildDepends      string
	buildDependsArch  string
	buildDependsIndep string
}

type SourcefileInfo struct {
	Path   string
	Size   int
	SHA256 string
}

func (sp SourcePackage) Files(prefix string) []SourcefileInfo {
	var ret []SourcefileInfo
	for _, f := range sp.files256 {
		fs := getArrayString(f, " ")
		if len(fs) != 3 {
			DebugPrintf("Unknown deb-src-files: %q", f)
			continue
		}
		s, _ := strconv.Atoi(fs[1])
		ret = append(ret, SourcefileInfo{
			Path:   path.Join(prefix, sp.directory, fs[2]),
			Size:   s,
			SHA256: fs[0],
		})
	}
	return ret
}

type Architecture string

type Architectures []string

func (as Architectures) Len() int {
	return len(as)
}
func (as Architectures) Swap(i, j int) {
	as[i], as[j] = as[j], as[i]
}

func (as Architectures) Less(i, j int) bool {
	return strings.Contains(string(as[i]), "64")
}

func (t BinaryPackage) valid() error {
	if len(t.Package) < 2 {
		return fmt.Errorf("W: pacakge name must be at least two characters long and start with an alphanumeric character: %q", t.Package)
	}
	if t.Filename == "" {
		return fmt.Errorf("W: BinaryPackage %q hasn't a filename field", t.Package)
	}
	if len(t.SHA256) != 64 {
		return fmt.Errorf("W: Wrong SHA256 length %s %d", t.Package, len(t.SHA256))
	}
	return nil
}

func (cf ControlFile) ToBinary() (BinaryPackage, error) {
	t := BinaryPackage{}
	t.Package = cf.Get("package")
	t.Version = cf.Get("version")
	t.Source, t.SourceVersion = parseSourceLine(cf.Get("source"), t.Package, t.Version)

	t.InstalledSize, _ = strconv.Atoi(cf.Get("installed-size"))
	t.Size, _ = strconv.Atoi(cf.Get("size"))

	for _, arch := range cf.GetArray("architecture", " ") {
		t.Architectures = append(t.Architectures, arch)
	}
	t.Description = cf.Get("description")
	t.Filename = cf.Get("filename")
	t.Tag = cf.Get("tag")
	t.Homepage = cf.Get("homepage")
	t.SHA256 = cf.Get("sha256")
	t.Maintainer = cf.Get("maintainer")
	t.Priority = cf.Get("priority")

	t.conflicts = cf.Get("conflicts")
	t.depends = cf.Get("depends")
	t.preDepends = cf.Get("pre-depends")
	return t, t.valid()
}
func parseProvides(str string) []string {
	if str == "" {
		return nil
	}

	depInfo, err := ParseDepInfo(str)
	if err != nil {
		panic(err)
	}
	var ret []string
	for depInfo.And != nil {
		ret = append(ret, depInfo.Name)
		depInfo = depInfo.And
	}
	ret = append(ret, depInfo.Name)
	return ret
}

func (cf SourcePackage) valid() error { return nil }
func (cf ControlFile) ToSource() (SourcePackage, error) {
	t := SourcePackage{}
	t.Package = cf.Get("package")
	t.Version = cf.Get("version")
	t.Description = cf.Get("description")
	t.Homepage = cf.Get("homepage")
	t.Format = cf.Get("format")
	t.Binary = cf.GetArray("binary", ",")
	if len(t.Binary) == 0 {
		t.Binary = []string{t.Package}
	}
	t.directory = cf.Get("directory")
	t.files256 = cf.GetMultiline("Checksums-Sha256")
	t.Architecture = cf.GetArray("architecture", " ")
	t.Maintainer = cf.Get("maintainer")
	t.Section = cf.Get("section")
	t.Priority = cf.Get("priority")
	t.buildDepends = cf.Get("build-depends")
	t.buildDependsArch = cf.Get("build-depends-arch")
	t.buildDependsIndep = cf.Get("build-depends-indep")

	plist := cf.GetMultiline("package-list")
	if len(plist) > 0 {
		for _, line := range plist {
			i, err := buildPackageListItem(line, t.Format)
			if err != nil {
				return t, FormatError{"SourcePackage", t.Package + t.Format, err}
			}
			if len(i.Archs) == 0 {
				i.Archs = t.Architecture
			}
			t.PackageList = append(t.PackageList, i)
		}
	} else {
		// from binary field to build
		for _, b := range t.Binary {
			t.PackageList = append(t.PackageList, PackageListItem{
				Name:  b,
				Ptype: "deb",
				Archs: t.Architecture,
			})
		}
	}
	return t, t.valid()
}

func (cf SourcePackage) GetBinary(arch string) []string {
	var ret []string
	for _, bp := range cf.PackageList {
		if bp.Ptype != "deb" {
			continue
		}
		if !bp.Support(arch) {
			continue
		}
		ret = append(ret, bp.Name)
	}
	if len(ret) == 0 {
		panic(fmt.Sprintf("GetBinary(%q) of %q failed. %d %+v\n", arch, cf.Package, len(cf.PackageList), cf.PackageList))
	}
	return ret
}

func (cf SourcePackage) BuildDepends(arch string, profile string) (*DepInfo, error) {
	AssertNoUseAny(arch)
	var deps []string
	switch arch {
	case "linux-all", "all":
		if len(cf.buildDependsIndep) > 0 {
			deps = append(deps, cf.buildDependsIndep)
		}
	default:
		if len(cf.buildDependsArch) > 0 {
			deps = append(deps, cf.buildDependsArch)
		}
	}
	var info *DepInfo
	var err error

	if len(deps) == 0 {
		info, err = ParseDepInfo(cf.buildDepends)
	} else {
		info, err = ParseDepInfo(strings.Join(deps, ","))
	}
	if err != nil || info == nil {
		return info, err
	}
	return info.Filter(arch, profile), err
}

func (cf BinaryPackage) deprecatedConflict() []string {
	info, _ := ParseDepInfo(cf.conflicts)
	if info == nil {
		return nil
	}
	return info.SimpleDeps()
}

func (cf BinaryPackage) Depends(arch string, profile string) (*DepInfo, error) {
	AssertNoUseAny(arch)
	var deps []string
	if cf.depends != "" {
		deps = append(deps, cf.depends)
	}
	if cf.preDepends != "" {
		deps = append(deps, cf.preDepends)
	}

	info, err := ParseDepInfo(strings.Join(deps, ","))
	if err != nil || info == nil {
		return info, err
	}

	return info.Filter(arch, profile), nil
}

func buildPackageListItem(line string, format string) (PackageListItem, error) {
	var r PackageListItem
	fields := getArrayString(line, " ")

	n := len(fields)
	if n < 4 || n > 7 {
		return r, FormatError{"PackageList", line, nil}
	}

	for i, v := range fields {
		switch i {
		case 0:
			r.Name = v
		case 1:
			r.Ptype = v
		case 2:
			r.Section = v
		case 3:
			r.Priority = v
		case 4:
			if !strings.HasPrefix(v, "arch=") {
				return r, FormatError{"PackageList", line, nil}
			}
			r.Archs = getArrayString(v[len("arch="):], ",")
		case 5:
			r.Profile = v
		case 6:
			r.Essional = v == "essential=yes"
		}
	}
	return r, nil
}
