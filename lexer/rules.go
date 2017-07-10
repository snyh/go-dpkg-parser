package lexer

import (
	"regexp"
	"sort"
	"strings"
)

type Token int

type Rules []Rule

type Rule struct {
	token           Token
	value           string
	caseInSensitive bool
	re              *regexp.Regexp
}

func (rs Rules) Match(raw string) (Token, int) {
	for _, r := range rs {
		width, ok := match(raw, r)
		if ok {
			return r.token, width
		}
	}
	return 0, 0
}

func (rs Rules) Len() int {
	return len(rs)
}

func (rs Rules) Less(i, j int) bool {
	return len(rs[i].value) > len(rs[j].value)

}
func (rs Rules) Swap(i, j int) {
	rs[j], rs[i] = rs[i], rs[j]
}

func (r Rules) Append(s []Rule) Rules {
	r = append(r, s...)
	sort.Sort(&r)
	return r
}

func match(raw string, rule Rule) (int, bool) {
	if rule.re != nil {
		return matchRegex(raw, rule.re)
	} else if rule.caseInSensitive {
		return matchI(raw, rule.value)
	} else {
		return matchC(raw, rule.value)
	}
}

func matchI(raw string, t string) (int, bool) {
	l := len(t)
	if l < len(raw) {
		return len(t),
			strings.HasPrefix(strings.ToLower(raw[:l]), strings.ToLower(t))
	} else {
		return len(t), strings.HasPrefix(strings.ToLower(raw), strings.ToLower(t))
	}
}

func matchC(raw string, t string) (int, bool) {
	return len(t), strings.HasPrefix(raw, t)
}

func matchRegex(raw string, re *regexp.Regexp) (int, bool) {
	r := re.FindStringIndex(raw)
	if len(r) == 0 {
		return 0, false
	}
	return r[0] + r[1], true
}

func (s *SimpleLexer) Add(token int, strs []string) {
	var rules []Rule
	for _, v := range strs {
		rules = append(rules, Rule{
			token: Token(token),
			value: v,
		})
	}
	s.rules = s.rules.Append(rules)
}

// AddI case insensitive match
func (s *SimpleLexer) AddI(token int, strs []string) {
	var rules []Rule
	for _, v := range strs {
		rules = append(rules, Rule{
			token:           Token(token),
			value:           v,
			caseInSensitive: true,
		})
	}
	s.rules = s.rules.Append(rules)
}

func (s *SimpleLexer) AddR(token int, res []string) {
	var rules []Rule
	for _, re := range res {
		rules = append(rules, Rule{
			token: Token(token),
			re:    regexp.MustCompilePOSIX(re),
		})
	}
	s.rules = s.rules.Append(rules)
}
