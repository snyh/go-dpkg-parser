package dpkg

import (
	"fmt"
	"github.com/bbuck/go-lexer"
	"strings"
)

type MM struct {
	*lexer.L
	info DepInfo
}

func lexDot(l *lexer.L) lexer.StateFunc {
	switch c := l.Peek(); c {
	case '!', '(', ')',
		'[', ']', '|':
		l.Emit(lexer.TokenType(c))
		return lexString
	default:
		return nil
	}
}

func lexString(l *lexer.L) lexer.StateFunc {
	digits := "0123456789-."
	alpha := "abcdefghijklmnopqrstuvwxyz"
	alpha2 := strings.ToUpper(alpha)
	l.Take(digits + alpha + alpha2)
	l.Emit(lexer.TokenType(ALPHA_NUMERIC))
	return nil
}

func parseDepInfo(str string) (DepInfo, error) {
	m := &MM{L: lexer.New(str, lexString)}
	m.Start()

	verNewParser().Parse(m)
	if m.Err != nil {
		return m.info, fmt.Errorf("Parsing %q failed: %v", str, m.Err)
	}
	return m.info, nil
}
func saveResult(l verLexer, r []Depend) {
	m := l.(*MM)
	m.info = DepInfo{
		Name: r[0].Name,
	}
}

func (m MM) Error(err string) {
	fmt.Println("ERR", err)
}

func (m MM) Lex(lval *verSymType) int {
	tok, done := m.NextToken()
	if done {
		return 0
	} else {
		lval.val = tok.Value
		return int(tok.Type)
	}
}

type DepInfo struct {
	Name    string
	VerMini string
	VerMax  string

	Arch    string
	Profile string
}

func (info DepInfo) String() string {
	return info.Name
}

func (info DepInfo) Match(arch string, profile string) bool {
	if info.Name != "" {
		return true
	}
	return false
}
