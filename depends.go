package dpkg

import (
	"fmt"
	"strings"
)

func (a Archive) CheckDep(str string) error {
	if str == "" {
		return nil
	}
	if v, ok := a.cache[str]; ok {
		return v
	}

	info, err := ParseDepInfo(str)
	if err != nil {
		return err
	}
	info = info.Filter(a.Architecture, "")
	err = a.checkDep(info)
	a.cache[str] = err
	return err
}
func (a Archive) hasPackage(name string) bool {
	_, ok := a.FindControl(name)
	if ok {
		return ok
	}
	return len(a.FindProvider(name)) != 0
}

func (a Archive) checkDep(info *DepInfo) error {
	if info == nil {
		return nil
	}

	str := info.String()
	if v, ok := a.cache[str]; ok {
		return v
	}

	if a.hasPackage(info.Name) {
		return a.record(str, a.checkDep(info.And))
	}

	if info.Or != nil {
		return a.record(str, a.checkDep(info.Or))
	} else {
		return a.record(str, fmt.Errorf("Can't find package %q", info.Name))
	}
}

func (a Archive) record(key string, v error) error {
	a.cache[key] = v
	return v
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
	Name string
	Ver  string
	Arch string

	Restrict struct {
		Archs    []string
		Profiles []string
	}

	And *DepInfo
	Or  *DepInfo
}

func (info DepInfo) SimpleDeps() []string {
	i := &info
	var ret []string
	for i != nil {
		ret = append(ret, i.Name)
		if i.And == nil {
			break
		}
		i = i.And
	}
	return ret
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

func (info DepInfo) Filter(arch string, profile string) *DepInfo {
	return depInfoFilter(&info, func(di *DepInfo) bool {
		return di != nil && di.match(arch, profile)
	})
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
	if info == nil {
		return nil
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
func (info DepInfo) matchProfile(profile string) bool {
	if len(info.Restrict.Profiles) == 0 {
		return true
	}
	for _, i := range info.Restrict.Profiles {
		if i == profile {
			return true
		}
	}
	return false
}
func (info DepInfo) matchArch(arch string) bool {
	if len(info.Restrict.Archs) == 0 {
		return true
	}
	for _, i := range info.Restrict.Archs {
		if i == arch {
			return true
		}
	}
	return false
}
