//go:generate go tool yacc -p ver ver.go.y
//
package dpkg

type Version struct {
	Minimal string
	Maximal string
	Arch    string
}

type DepInfo struct {
	Name    string
	VerMini string
	VerMax  string

	Arch    string
	Profile string
}

func (DepInfo) Match(arch string, profile string) bool {
	return false
}

func ParseDepInfo(str string) (DepInfo, error) {
	panic("Not Implement")
}
