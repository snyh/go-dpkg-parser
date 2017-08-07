package dpkg

import (
	"fmt"
	"strings"
)

type DependStatus struct {
	satisfy bool
	err     error
	chain   *DependStatus

	result []struct {
		Name    string
		Version string
	}
}

func (a Archive) parseDepend(name string) error {
	return nil
}

func AssertNoUseAny(arch string) {
	if arch == "any" {
		panic("It's wrong to query depends by architecture of any.")
	}
}

func parseSourceLine(str string, defSource, defVer string) (string, string) {
	// TODO: re implement by parseDepInfo
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

type DepInfo struct {
	Name   string
	VerMin string
	VerMax string

	Archs    []string
	Profiles []string

	And *DepInfo
	Or  *DepInfo
}

func (info DepInfo) String() string {
	r := info.Name
	if info.Or != nil {
		r += " | " + info.Or.String()
	}
	if info.And != nil {
		r += ", " + info.And.String()
	}
	return r
}
func (info DepInfo) matchProfile(profile string) bool {
	if len(info.Profiles) == 0 {
		return true
	}
	for _, i := range info.Profiles {
		if i == profile {
			return true
		}
	}
	return false
}
func (info DepInfo) matchArch(arch string) bool {
	if len(info.Archs) == 0 {
		return true
	}
	for _, i := range info.Archs {
		if i == arch {
			return true
		}
	}
	return false
}

func (info DepInfo) Filter(arch string, profile string) (DepInfo, error) {
	r := depInfoFilter(&info, func(di *DepInfo) bool {
		return di != nil && di.match(arch, profile)
	})
	if r != nil {
		return *r, nil
	}
	return info, fmt.Errorf("Empty Result")
}

type filterFunc func(*DepInfo) bool

func depInfoFilter(info *DepInfo, fn filterFunc) *DepInfo {
	if info == nil {
		return info
	}
	info = depInfoFilterOr(info, fn)
	if fn(info) {
		info.And = depInfoFilter(info.And, fn)
		return info
	}
	return depInfoFilter(info.And, fn)
}

func depInfoFilterOr(info *DepInfo, fn filterFunc) *DepInfo {
	if info == nil {
		return info
	}
	if fn(info) {
		info.Or = depInfoFilterOr(info.Or, fn)
		return info
	}

	if info.Or == nil {
		return depInfoFilterOr(info.And, fn)
	}

	info.Or.And = info.And
	return depInfoFilterOr(info.Or, fn)
}

func (info DepInfo) match(arch string, profile string) bool {
	return info.Name != "" && info.matchArch(arch) && info.matchProfile(profile)
}
