//go:generate go tool yacc -p ver ver.go.y
//
package dpkg

import (
	lexer "./lexer"
	"fmt"
)

func ParseDepInfo(str string) (DepInfo, error) {
	deps, err := parseDepends(str)
	if err != nil {
		return DepInfo{}, err
	}
	return DepInfo{
		Name: deps[0].Name,
	}, nil
}

type myLexer struct {
	e    string
	slex *lexer.SimpleLexer
	r    []Depend
}

func NewMyLexer(raw string, archList, profileList []string) *myLexer {
	sl := lexer.NewSimpleLexer("dsc-version")
	sl.SetInputString(raw)

	sl.AddIgnores(" ")

	sl.Add(ANY, []string{"any"})
	sl.Add(NATIVE, []string{"native"})

	sl.SetBasicToken("()[]<>=!,:+-~.")

	sl.AddI(ARCH_NAME, archList)
	sl.AddI(PROFILE, profileList)
	sl.AddR(ALPHA_NUMERIC, []string{`^[0-9a-zA-Z]+`})
	return &myLexer{
		slex: sl,
	}
}

func (l *myLexer) Lex(lval *verSymType) int {
	t, s := l.slex.Token()
	lval.val = s
	return int(t)
}

func (l *myLexer) Error(err string) {
	l.e = err
}

type Depend struct {
	Name      string
	Version   string
	Operation string
}
type Version struct {
	Minimal string
	Maximal string
	Arch    string
}

func parseDepends(v string) ([]Depend, error) {
	archList := []string{"amd64", "i386", "ia64",
		"mips", "mipsel", "mips64el", "mips64", "mipsn32", "mipsn32el",
		"kfreebsd-amd64", "kfreebsd-i386",
		"s390", "s390x", "sparc",
		"arm", "armel", "armhf", "arm64",
		"powerpc", "powerpcspe",
		"ppc64el", "ppc64", "ppc64el",
		"linux-any", "hurd-any", "kfreebsd-any", "any-i386", "any-amd64",
		"hurd-i386", "x32", "sparc64",
		"m68k", "sh4", "alpha", "hppa", "avr32",
		"or1k",
	}
	profileList := []string{"cross", "stage1", "nocheck", "stage2", "nodoc", "profile.nobluetooth", "nobiarch"}
	l := NewMyLexer(v, archList, profileList)
	verNewParser().Parse(l)
	if l.e != "" {
		return nil, fmt.Errorf("Parsing %q failed: %v", v, l.e)
	}
	return l.r, nil
}
