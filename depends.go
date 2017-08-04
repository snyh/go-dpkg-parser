package dpkg

import (
	"strings"
)

func AssertNoUseAny(arch string) {
	if arch == "any" {
		panic("It's wrong to query depends by architecture of any.")
	}
}

func matchDepends(rawDeps []string, arch string, profile string) ([]DepInfo, error) {
	var ret []DepInfo
	for _, raw := range rawDeps {
		info, err := parseDepInfo(raw)
		if err != nil {
			return nil, err
		}
		if info.Match(arch, profile) {
			ret = append(ret, info)
		}
	}
	return ret, nil
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
}

func (info DepInfo) String() string {
	return info.Name
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

func (info DepInfo) Match(arch string, profile string) bool {
	return info.Name != "" && info.matchArch(arch) && info.matchProfile(profile)
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
