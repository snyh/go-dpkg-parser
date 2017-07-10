//go:generate go tool yacc -o parser_error.go -p error parser_error.y
package lexer

import (
	"fmt"
)

type ruleFlag int

const (
	ruleFlagConst ruleFlag = iota
	ruleFlagRegex
	ruleFlagCaseInsensitive
)

type RecordState struct {
	Offset int
	Width  int
	Token  Token
}

func (s RecordState) Str(raw string) string {
	if len(raw) < s.Offset+s.Width {
		return ""
	}
	return raw[s.Offset : s.Offset+s.Width]
}

func (s RecordState) Remain(raw string) string {
	return raw[s.Offset+s.Width:]
}

type SimpleLexer struct {
	name string

	operators map[rune]Token
	ignores   map[rune]struct{}
	rules     Rules

	Input   string
	Pos     int
	Records []RecordState

	Debug bool
}

func NewSimpleLexer(name string) *SimpleLexer {
	s := &SimpleLexer{
		name:      name,
		operators: make(map[rune]Token),
		ignores:   make(map[rune]struct{}),
	}
	return s
}

func (s *SimpleLexer) SetInputString(r string) {
	s.Input = r
	s.Pos = 0
	s.Records = nil
}

func (s *SimpleLexer) SetBasicToken(tokens string) {
	for _, t := range tokens {
		s.operators[t] = Token(t)
	}
}

func (s *SimpleLexer) AddIgnores(bs string) {
	for _, ch := range bs {
		s.ignores[ch] = struct{}{}
	}
}

func (s *SimpleLexer) Token() (Token, string) {
	s.doIgnore()

	if s.Pos >= len(s.Input) {
		return 0, ""
	}

	if t, ok := s.operators[rune(s.Input[s.Pos])]; ok {
		if s.Debug {
			fmt.Printf("Hit operator ...%q\n", string([]byte{s.Input[s.Pos]}))
		}
		if ok {
			return s.recordAndReturn(t, 1)
		}
	} else {
		if s.Debug {
			fmt.Printf("Not hit operator ...%q\n", string([]byte{s.Input[s.Pos]}))
		}
	}
	return s.recordAndReturn(s.rules.Match(s.Input[s.Pos:]))
}

func (l *SimpleLexer) recordAndReturn(token Token, width int) (Token, string) {
	if token == 0 || width == 0 {
		return 0, ""
	}
	s := RecordState{
		Token:  token,
		Offset: l.Pos,
		Width:  width,
	}

	l.Records = append(l.Records, s)

	l.Pos += width

	return token, s.Str(l.Input)
}

func (l *SimpleLexer) doIgnore() bool {
	if len(l.ignores) == 0 {
		return false
	}
	before := l.Pos
	for i := l.Pos; i < len(l.Input); i++ {
		ch := l.Input[i]
		if _, ok := l.ignores[rune(ch)]; !ok {
			break
		}

		l.Pos++
	}
	return before != l.Pos
}
