//go:generate go tool yacc -p ver ver.go.y
//
package dpkg

type Version struct {
	Minimal string
	Maximal string
	Arch    string
}

func FindPkg(name string, verS string, verE string, arch string) {
}

func (db *PackageDatabase) ListAllVersion() []string {
	var ret []string
	for _, cf := range db.SourcePackages {
		ret = append(ret, cf.GetString("version"))
	}
	return ret
}
