package dpkg

import (
	"fmt"
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
	for _, i := range item.Archs {
		if i == "any" || i == "all" {
			return true
		}
		if i == string(arch) {
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
	Architecture string            `json:"architecture"`
	Maintainer   string            `json:"maintainer"`
	Format       string            `json:"format"`
	Binary       []string          `json:"binary"`
	PackageList  []PackageListItem `json:"package_list"`

	Section  string `json:"section"`
	Priority string `json:"priority"`

	buildDepends []string
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

func parseSourceLine(str string, defSource, defVer string) (string, string) {
	if str == "" {
		return defSource, defVer
	}
	fs := getArrayString(str, " ")
	switch len(fs) {
	case 2:
		return fs[0], strings.Trim(fs[1], "()")
	case 1:
		return fs[0], defVer
	default:
		DebugPrintf("Invalid source line %q (%d)\n", str, len(fs))
		return defSource, defVer
	}
}

func (cf ControlFile) ToBinary() (BinaryPackage, error) {
	t := BinaryPackage{}
	t.Package = cf.GetString("package")
	t.Version = cf.GetString("version")
	t.Source, t.SourceVersion = parseSourceLine(cf.GetString("source"), t.Package, t.Version)

	t.InstalledSize, _ = strconv.Atoi(cf.GetString("installed-size"))
	t.Size, _ = strconv.Atoi(cf.GetString("size"))

	for _, arch := range cf.GetArrayString("architecture", " ") {
		t.Architectures = append(t.Architectures, arch)
	}
	t.Description = cf.GetString("description")
	t.Filename = cf.GetString("filename")
	t.Tag = cf.GetString("tag")
	t.Homepage = cf.GetString("homepage")
	t.SHA256 = cf.GetString("sha256")
	t.Maintainer = cf.GetString("maintainer")

	return t, t.valid()
}

func (cf SourcePackage) valid() error { return nil }
func (cf ControlFile) ToSource() (SourcePackage, error) {
	t := SourcePackage{}
	t.Package = cf.GetString("package")
	t.Version = cf.GetString("version")
	t.Description = cf.GetString("description")
	t.Homepage = cf.GetString("homepage")
	t.Format = cf.GetString("format")
	t.Binary = cf.GetArrayString("binary", ",")
	t.Architecture = cf.GetString("architecture")
	t.Maintainer = cf.GetString("maintainer")
	t.Section = cf.GetString("section")
	t.Priority = cf.GetString("priority")
	t.buildDepends = cf.GetArrayString("build-depends", ",")

	plist := cf.GetMultiline("package-list")
	if len(plist) > 0 {
		for _, line := range plist {
			i, err := buildPackageListItem(line, t.Format)
			if err != nil {
				return t, FormatError{"SourcePackage", t.Package + t.Format, err}
			}
			t.PackageList = append(t.PackageList, i)
		}
	} else {
		// from binary field to build
		for _, b := range t.Binary {
			t.PackageList = append(t.PackageList, PackageListItem{
				Name:  b,
				Archs: []string{t.Architecture},
			})
		}
	}

	return t, t.valid()
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

func (cf SourcePackage) GetBinary(arch string) []string {
	var ret []string
	// arch <- [ "amd64", "i386" ]
	for _, bp := range cf.PackageList {
		if bp.Ptype != "deb" {
			continue
		}
		if !bp.Support(arch) {
			continue
		}
		ret = append(ret, bp.Name)
	}
	return ret
}

func (cf SourcePackage) BuildDepends(arch string, profile string) ([]DepInfo, error) {
	var ret []DepInfo
	for _, raw := range cf.buildDepends {
		info, err := ParseDepInfo(raw)
		if err != nil {
			return nil, err
		}
		if info.Match(arch, profile) {
			ret = append(ret, info)
		}
	}
	return ret, nil
}
