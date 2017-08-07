//go:generate go tool yacc -p ver -o verparser.go ver.go.y
//
package dpkg

import (
	"fmt"
	"github.com/bbuck/go-lexer"
	"strings"
)

type MM struct {
	*lexer.L
	info DepInfo
	str  string
}

const (
	_DIGIT        = "0123456789"
	_ALPHA        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	_ALPHANUM     = _DIGIT + _ALPHA
	_PKGNAMESPECS = "+-~."
	_PKGNAME      = _ALPHANUM + _PKGNAMESPECS
)

func lexPkgName(l *lexer.L) lexer.StateFunc {
	ignoreSpaces(l)
	if !takeAny(l, _PKGNAME) {
		return nil
	}
	l.Emit(lexer.TokenType(PKGNAME))
	return lexPkgOthers
}

func lexPkgOthers(l *lexer.L) lexer.StateFunc {
	ignoreSpaces(l)
	switch l.Peek() {
	case '[':
		if takeMatch(l, ARCH_QUALIFIER, '[', ']') {
			return lexPkgOthers
		}
	case '(':
		if takeMatch(l, VERSION, '(', ')') {
			return lexPkgOthers
		}
	case '<':
		if takeMatch(l, PROFILE, '<', '>') {
			return lexPkgOthers
		}
	case '|':
		l.Next()
		l.Emit(lexer.TokenType('|'))
		return lexPkgName
	}
	return nil
}

func parseDepInfo(str string) (DepInfo, error) {
	m := &MM{
		L:   lexer.New(str, lexPkgName),
		str: str,
	}
	m.Start()

	verNewParser().Parse(m)
	if m.Err != nil {
		return m.info, fmt.Errorf("Parsing %q failed: %v", str, m.Err)
	}
	return m.info, nil
}

func saveResult(l verLexer, info DepInfo) {
	m := l.(*MM)
	m.info = info
}

func (m MM) Error(err string) {
	fmt.Println("ERR", err)
}

func (m MM) Lex(lval *verSymType) int {
	tok, done := m.NextToken()
	//	fmt.Printf("%q TOKEN: %v %v\n", m.str, tok, done)
	if done {
		return 0
	} else {
		lval.val = tok.Value
		return int(tok.Type)
	}
}

func ignoreSpaces(l *lexer.L) {
	l.Take(" ")
	l.Ignore()
}

func takeAny(l *lexer.L, strs string) bool {
	r := l.Next()
	found := false
	for strings.ContainsRune(strs, r) {
		r = l.Next()
		found = true
	}
	l.Rewind()
	return found
}

func takeMatch(l *lexer.L, tok int, left rune, right rune) bool {
	c := l.Peek()
	if c != left {
		return false
	}
	l.Next()
	l.Ignore()
	for i := 0; i < 1000; i++ {
		if l.Peek() == right {
			l.Emit(lexer.TokenType(tok))
			l.Next()
			l.Ignore()
			return true
		}
		l.Next()
	}
	return false
}
