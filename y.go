//line ver.go.y:2
package dpkg

import __yyfmt__ "fmt"

//line ver.go.y:2
import "./lexer"

//line ver.go.y:19
type verSymType struct {
	yys int
	val string
	r   []Depend
	tmp Depend
}

const ANY = 57346
const NATIVE = 57347
const VER_NUM = 57348
const PKG_NAME = 57349
const PROFILE = 57350
const ARCH_NAME = 57351

var verToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"ANY",
	"NATIVE",
	"','",
	"'!'",
	"')'",
	"'('",
	"']'",
	"'['",
	"'<'",
	"'>'",
	"VER_NUM",
	"PKG_NAME",
	"PROFILE",
	"ARCH_NAME",
	"'|'",
	"':'",
	"'='",
}
var verStatenames = [...]string{}

const verEofCode = 1
const verErrCode = 2
const verInitialStackSize = 16

//line ver.go.y:110

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

	sl.SetBasicToken("()[]<>=!,")

	sl.AddR(VER_NUM, []string{})
	sl.AddI(ARCH_NAME, archList)
	sl.AddI(PROFILE, profileList)
	sl.AddR(PKG_NAME, []string{`^[a-z][a-z0-9\.\-\+]+`})
	sl.AddR(VER_NUM, []string{`^[0-9]+[a-z0-9\.\+\-\~\:]*`})
	return &myLexer{
		slex: sl,
	}
}

func saveResult(l verLexer, r []Depend) {
	l.(*myLexer).r = r
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

func ParseDepends(v string) ([]Depend, error) {
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
		return nil, __yyfmt__.Errorf("Parsing %q failed: %v", v, l.e)
	}
	return l.r, nil
}

//line yacctab:1
var verExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const verNprod = 23
const verPrivate = 57344

var verTokenNames []string
var verStates []string

const verLast = 40

var verAct = [...]int{

	21, 8, 28, 9, 10, 30, 14, 15, 23, 29,
	6, 7, 26, 31, 36, 20, 19, 18, 22, 16,
	40, 25, 4, 33, 34, 39, 38, 37, 35, 32,
	27, 2, 5, 1, 24, 17, 13, 3, 12, 11,
}
var verPact = [...]int{

	7, -1000, 26, -8, -1000, 7, 7, 2, 3, 1,
	5, -1000, -1000, -1000, -1000, -1000, -1000, 22, -1000, -11,
	-7, 19, 1, 1, 15, -1000, -2, -1000, 13, 12,
	11, 6, -1000, -1000, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000,
}
var verPgo = [...]int{

	0, 33, 31, 37, 36, 35, 0, 34,
}
var verR1 = [...]int{

	0, 1, 1, 2, 2, 3, 3, 3, 3, 3,
	4, 4, 4, 5, 5, 5, 5, 5, 6, 6,
	6, 7, 7,
}
var verR2 = [...]int{

	0, 3, 1, 3, 1, 1, 3, 4, 4, 4,
	1, 1, 1, 1, 3, 3, 3, 3, 1, 2,
	2, 1, 2,
}
var verChk = [...]int{

	-1000, -1, -2, -3, 15, 6, 18, 19, 9, 11,
	12, -1, -2, -4, 4, 5, 17, -5, 14, 13,
	12, -6, 17, 7, -7, 16, 7, 8, 13, 20,
	12, 20, 10, -6, -6, 13, 16, 14, 14, 14,
	14,
}
var verDef = [...]int{

	0, -2, 2, 4, 5, 0, 0, 0, 0, 0,
	0, 1, 3, 6, 10, 11, 12, 0, 13, 0,
	0, 0, 18, 0, 0, 21, 0, 7, 0, 0,
	0, 0, 8, 20, 19, 9, 22, 14, 15, 16,
	17,
}
var verTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 7, 3, 3, 3, 3, 3, 3,
	9, 8, 3, 3, 6, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 19, 3,
	12, 20, 13, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 11, 3, 10, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 18,
}
var verTok2 = [...]int{

	2, 3, 4, 5, 14, 15, 16, 17,
}
var verTok3 = [...]int{
	0,
}

var verErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	verDebug        = 0
	verErrorVerbose = false
)

type verLexer interface {
	Lex(lval *verSymType) int
	Error(s string)
}

type verParser interface {
	Parse(verLexer) int
	Lookahead() int
}

type verParserImpl struct {
	lval  verSymType
	stack [verInitialStackSize]verSymType
	char  int
}

func (p *verParserImpl) Lookahead() int {
	return p.char
}

func verNewParser() verParser {
	return &verParserImpl{}
}

const verFlag = -1000

func verTokname(c int) string {
	if c >= 1 && c-1 < len(verToknames) {
		if verToknames[c-1] != "" {
			return verToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func verStatname(s int) string {
	if s >= 0 && s < len(verStatenames) {
		if verStatenames[s] != "" {
			return verStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func verErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !verErrorVerbose {
		return "syntax error"
	}

	for _, e := range verErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + verTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := verPact[state]
	for tok := TOKSTART; tok-1 < len(verToknames); tok++ {
		if n := base + tok; n >= 0 && n < verLast && verChk[verAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if verDef[state] == -2 {
		i := 0
		for verExca[i] != -1 || verExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; verExca[i] >= 0; i += 2 {
			tok := verExca[i]
			if tok < TOKSTART || verExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if verExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += verTokname(tok)
	}
	return res
}

func verlex1(lex verLexer, lval *verSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = verTok1[0]
		goto out
	}
	if char < len(verTok1) {
		token = verTok1[char]
		goto out
	}
	if char >= verPrivate {
		if char < verPrivate+len(verTok2) {
			token = verTok2[char-verPrivate]
			goto out
		}
	}
	for i := 0; i < len(verTok3); i += 2 {
		token = verTok3[i+0]
		if token == char {
			token = verTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = verTok2[1] /* unknown char */
	}
	if verDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", verTokname(token), uint(char))
	}
	return char, token
}

func verParse(verlex verLexer) int {
	return verNewParser().Parse(verlex)
}

func (verrcvr *verParserImpl) Parse(verlex verLexer) int {
	var vern int
	var verVAL verSymType
	var verDollar []verSymType
	_ = verDollar // silence set and not used
	verS := verrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	verstate := 0
	verrcvr.char = -1
	vertoken := -1 // verrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		verstate = -1
		verrcvr.char = -1
		vertoken = -1
	}()
	verp := -1
	goto verstack

ret0:
	return 0

ret1:
	return 1

verstack:
	/* put a state and value onto the stack */
	if verDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", verTokname(vertoken), verStatname(verstate))
	}

	verp++
	if verp >= len(verS) {
		nyys := make([]verSymType, len(verS)*2)
		copy(nyys, verS)
		verS = nyys
	}
	verS[verp] = verVAL
	verS[verp].yys = verstate

vernewstate:
	vern = verPact[verstate]
	if vern <= verFlag {
		goto verdefault /* simple state */
	}
	if verrcvr.char < 0 {
		verrcvr.char, vertoken = verlex1(verlex, &verrcvr.lval)
	}
	vern += vertoken
	if vern < 0 || vern >= verLast {
		goto verdefault
	}
	vern = verAct[vern]
	if verChk[vern] == vertoken { /* valid shift */
		verrcvr.char = -1
		vertoken = -1
		verVAL = verrcvr.lval
		verstate = vern
		if Errflag > 0 {
			Errflag--
		}
		goto verstack
	}

verdefault:
	/* default state action */
	vern = verDef[verstate]
	if vern == -2 {
		if verrcvr.char < 0 {
			verrcvr.char, vertoken = verlex1(verlex, &verrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if verExca[xi+0] == -1 && verExca[xi+1] == verstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			vern = verExca[xi+0]
			if vern < 0 || vern == vertoken {
				break
			}
		}
		vern = verExca[xi+1]
		if vern < 0 {
			goto ret0
		}
	}
	if vern == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			verlex.Error(verErrorMessage(verstate, vertoken))
			Nerrs++
			if verDebug >= 1 {
				__yyfmt__.Printf("%s", verStatname(verstate))
				__yyfmt__.Printf(" saw %s\n", verTokname(vertoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for verp >= 0 {
				vern = verPact[verS[verp].yys] + verErrCode
				if vern >= 0 && vern < verLast {
					verstate = verAct[vern] /* simulate a shift of "error" */
					if verChk[verstate] == verErrCode {
						goto verstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if verDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", verS[verp].yys)
				}
				verp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if verDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", verTokname(vertoken))
			}
			if vertoken == verEofCode {
				goto ret1
			}
			verrcvr.char = -1
			vertoken = -1
			goto vernewstate /* try again in the same state */
		}
	}

	/* reduction by production vern */
	if verDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", vern, verStatname(verstate))
	}

	vernt := vern
	verpt := verp
	_ = verpt // guard against "declared and not used"

	verp -= verR2[vern]
	// verp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if verp+1 >= len(verS) {
		nyys := make([]verSymType, len(verS)*2)
		copy(nyys, verS)
		verS = nyys
	}
	verVAL = verS[verp+1]

	/* consult goto table to find next state */
	vern = verR1[vern]
	verg := verPgo[vern]
	verj := verg + verS[verp].yys + 1

	if verj >= verLast {
		verstate = verAct[verg]
	} else {
		verstate = verAct[verj]
		if verChk[verstate] != -vern {
			verstate = verAct[verg]
		}
	}
	// dummy call; replaced with literal code
	switch vernt {

	case 1:
		verDollar = verS[verpt-3 : verpt+1]
		//line ver.go.y:31
		{
			verVAL.r = append([]Depend{verDollar[1].tmp}, verDollar[3].r...)
			saveResult(verlex, verVAL.r)
		}
	case 2:
		verDollar = verS[verpt-1 : verpt+1]
		//line ver.go.y:36
		{
			verVAL.r = []Depend{verDollar[1].tmp}
			saveResult(verlex, verVAL.r)
		}
	case 3:
		verDollar = verS[verpt-3 : verpt+1]
		//line ver.go.y:44
		{
			verVAL.r = append([]Depend{verDollar[1].tmp}, verDollar[3].r...)
		}
	case 4:
		verDollar = verS[verpt-1 : verpt+1]
		//line ver.go.y:48
		{
			verVAL.r = []Depend{verDollar[1].tmp}
		}
	case 5:
		verDollar = verS[verpt-1 : verpt+1]
		//line ver.go.y:54
		{
			verVAL.tmp.Name = verDollar[1].val
			verVAL.tmp.Version = ""
			verVAL.tmp.Operation = ""
		}
	case 7:
		verDollar = verS[verpt-4 : verpt+1]
		//line ver.go.y:61
		{
			verVAL.tmp.Name = verDollar[1].val
			verVAL.tmp.Version = verDollar[3].tmp.Version
			verVAL.tmp.Operation = verDollar[3].tmp.Operation
		}
	case 13:
		verDollar = verS[verpt-1 : verpt+1]
		//line ver.go.y:76
		{
			verVAL.tmp.Version = verDollar[1].val
		}
	case 14:
		verDollar = verS[verpt-3 : verpt+1]
		//line ver.go.y:80
		{
			verVAL.tmp.Version = verDollar[3].val
			verVAL.tmp.Operation = "GT"
		}
	case 15:
		verDollar = verS[verpt-3 : verpt+1]
		//line ver.go.y:85
		{
			verVAL.tmp.Version = verDollar[3].val
			verVAL.tmp.Operation = "GTE"
		}
	case 16:
		verDollar = verS[verpt-3 : verpt+1]
		//line ver.go.y:90
		{
			verVAL.tmp.Version = verDollar[3].val
			verVAL.tmp.Operation = "ST"
		}
	case 17:
		verDollar = verS[verpt-3 : verpt+1]
		//line ver.go.y:95
		{
			verVAL.tmp.Version = verDollar[3].val
			verVAL.tmp.Operation = "STE"
		}
	}
	goto verstack /* stack new state and value */
}
