package dpkg

import (
	"fmt"
	"strconv"
	"strings"
)

type BinaryPackage struct {
	Package       string        `json:"package"`
	Version       string        `json:"version"`
	InstalledSize int           `json:"installed_size"`
	Size          int           `json:"size"`
	Architectures Architectures `json:"architectures"`
	Description   string        `json:"description"`
	Filename      string        `json:"filename"`
	Tag           string        `json:"tag"`
	Homepage      string        `json:"homepage"`
}

type PackageListItem struct {
	Name     string
	Ptype    string
	Section  string
	Priority string
	Arch     string
}

type SourcePackage struct {
	Package     string `json:"package"`
	Version     string `json:"version"`
	Description string `json:"description"`
	Homepage    string `json:"homepage"`

	Format     string            `json:"format"`
	Binary     []string          `json:"binary"`
	PackgeList []PackageListItem `json:"package_list"`
}

type Architecture string

type Architectures []Architecture

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
	return nil
}

func (cf ControlFile) ToBinary() (BinaryPackage, error) {
	t := BinaryPackage{}
	t.Package = cf.GetString("package")
	t.Version = cf.GetString("version")
	t.InstalledSize, _ = strconv.Atoi(cf.GetString("installed-size"))
	t.Size, _ = strconv.Atoi(cf.GetString("size"))

	for _, arch := range cf.GetArrayString("architecture", " ") {
		t.Architectures = append(t.Architectures, Architecture(arch))
	}
	t.Description = cf.GetString("description")
	t.Filename = cf.GetString("filename")
	t.Tag = cf.GetString("tag")
	t.Homepage = cf.GetString("homepage")

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
	return t, t.valid()
}
