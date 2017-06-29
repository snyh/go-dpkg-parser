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
		return fmt.Errorf("W: parsing ControlFile not enough fields: %q %q %q", t.Package, t.Tag, t.Filename)
	}
	return nil
}

func (cf ControlFile) ToBinary() (BinaryPackage, error) {
	t := BinaryPackage{}
	t.Package = cf.GetString("package")
	t.Version = cf.GetString("version")
	t.InstalledSize, _ = strconv.Atoi(cf.GetString("installed-size"))
	t.Size, _ = strconv.Atoi(cf.GetString("size"))

	for _, arch := range cf.GetArrayString("architecture") {
		t.Architectures = append(t.Architectures, Architecture(arch))
	}
	t.Description = cf.GetString("description")
	t.Filename = cf.GetString("filename")
	t.Tag = cf.GetString("tag")
	t.Homepage = cf.GetString("homepage")

	return t, t.valid()
}
