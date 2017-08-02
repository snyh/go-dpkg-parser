//line ver.go.y:2
package dpkg

import __yyfmt__ "fmt"

//line ver.go.y:2
//line ver.go.y:19
type verSymType struct {
	yys int
	val string
	r   []Depend
	tmp Depend
}

const ANY = 57346
const NATIVE = 57347
const ALPHA_NUMERIC = 57348
const PROFILE = 57349
const ARCH_NAME = 57350

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
	"ALPHA_NUMERIC",
	"PROFILE",
	"ARCH_NAME",
	"'.'",
	"'+'",
	"'~'",
	"'|'",
	"'-'",
	"':'",
	"'='",
}
var verStatenames = [...]string{}

const verEofCode = 1
const verErrCode = 2
const verInitialStackSize = 16

//line ver.go.y:135

//line yacctab:1
var verExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
	-1, 23,
	21, 27,
	-2, 21,
	-1, 58,
	21, 23,
	-2, 29,
}

const verNprod = 38
const verPrivate = 57344

var verTokenNames []string
var verStates []string

const verLast = 70

var verAct = [...]int{

	56, 20, 26, 38, 23, 24, 8, 60, 9, 10,
	43, 36, 11, 13, 39, 48, 12, 6, 65, 7,
	28, 37, 62, 61, 63, 41, 40, 42, 31, 27,
	45, 46, 16, 17, 25, 47, 30, 58, 49, 50,
	51, 52, 44, 34, 18, 53, 54, 55, 59, 57,
	22, 21, 25, 33, 32, 5, 35, 29, 19, 15,
	2, 4, 64, 66, 67, 59, 57, 14, 3, 1,
}
var verPact = [...]int{

	41, -1000, -1000, -3, -5, -1000, 41, 28, 38, 13,
	21, 40, 39, 29, -1000, -1000, -1000, -1000, -1000, 48,
	-1000, -2, -9, 8, -11, -1000, 32, 13, 13, 22,
	-1000, 0, -1000, -1000, -1000, -1000, 20, 20, 20, 20,
	20, 20, 20, 23, -1000, -1000, -1000, -1000, -1000, -1000,
	-1000, -1000, -1000, -1000, -1000, -1000, -1000, -14, 5, 8,
	20, 4, 4, 4, -1000, 5, -1000, -1000,
}
var verPgo = [...]int{

	0, 69, 60, 68, 61, 59, 58, 2, 57, 1,
	4, 5, 0,
}
var verR1 = [...]int{

	0, 1, 2, 2, 4, 4, 4, 4, 3, 3,
	3, 3, 3, 5, 5, 5, 6, 6, 6, 6,
	6, 9, 9, 10, 10, 10, 10, 11, 11, 12,
	12, 12, 12, 7, 7, 7, 8, 8,
}
var verR2 = [...]int{

	0, 1, 3, 1, 1, 3, 3, 3, 1, 3,
	4, 4, 4, 1, 1, 1, 1, 3, 3, 3,
	3, 1, 3, 1, 3, 3, 3, 1, 3, 1,
	3, 3, 3, 1, 2, 2, 1, 2,
}
var verChk = [...]int{

	-1000, -1, -2, -3, -4, 14, 20, 22, 9, 11,
	12, 17, 21, 18, -2, -5, 4, 5, 16, -6,
	-9, 13, 12, -10, -11, 14, -7, 16, 7, -8,
	15, 7, 14, 14, 14, 8, 13, 23, 12, 23,
	18, 17, 19, 21, 10, -7, -7, 13, 15, -9,
	-9, -9, -9, -10, -10, -10, -12, -11, 14, -10,
	21, 18, 17, 19, -12, 14, -12, -12,
}
var verDef = [...]int{

	0, -2, 1, 3, 8, 4, 0, 0, 0, 0,
	0, 0, 0, 0, 2, 9, 13, 14, 15, 0,
	16, 0, 0, -2, 0, 23, 0, 33, 0, 0,
	36, 0, 5, 6, 7, 10, 0, 0, 0, 0,
	0, 0, 0, 0, 11, 35, 34, 12, 37, 17,
	18, 19, 20, 24, 25, 26, 22, 0, -2, 27,
	0, 0, 0, 0, 30, 29, 31, 32,
}
var verTok1 = [...]int{

	1, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 7, 3, 3, 3, 3, 3, 3,
	9, 8, 3, 18, 6, 21, 17, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 22, 3,
	12, 23, 13, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 11, 3, 10, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	3, 3, 3, 3, 20, 3, 19,
}
var verTok2 = [...]int{

	2, 3, 4, 5, 14, 15, 16,
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
		verDollar = verS[verpt-1 : verpt+1]
		//line ver.go.y:31
		{
			verVAL.r = []Depend{verDollar[1].tmp}
			saveResult(verlex, verVAL.r)
		}
	case 2:
		verDollar = verS[verpt-3 : verpt+1]
		//line ver.go.y:39
		{
			verVAL.r = append([]Depend{verDollar[1].tmp}, verDollar[3].r...)
		}
	case 3:
		verDollar = verS[verpt-1 : verpt+1]
		//line ver.go.y:43
		{
			verVAL.r = []Depend{verDollar[1].tmp}
		}
	case 8:
		verDollar = verS[verpt-1 : verpt+1]
		//line ver.go.y:55
		{
			verVAL.tmp.Name = verDollar[1].val
			verVAL.tmp.Version = ""
			verVAL.tmp.Operation = ""
		}
	case 10:
		verDollar = verS[verpt-4 : verpt+1]
		//line ver.go.y:62
		{
			verVAL.tmp.Name = verDollar[1].val
			verVAL.tmp.Version = verDollar[3].tmp.Version
			verVAL.tmp.Operation = verDollar[3].tmp.Operation
		}
	case 16:
		verDollar = verS[verpt-1 : verpt+1]
		//line ver.go.y:77
		{
			verVAL.tmp.Version = verDollar[1].val
		}
	case 17:
		verDollar = verS[verpt-3 : verpt+1]
		//line ver.go.y:81
		{
			verVAL.tmp.Version = verDollar[3].val
			verVAL.tmp.Operation = "GT"
		}
	case 18:
		verDollar = verS[verpt-3 : verpt+1]
		//line ver.go.y:86
		{
			verVAL.tmp.Version = verDollar[3].val
			verVAL.tmp.Operation = "GTE"
		}
	case 19:
		verDollar = verS[verpt-3 : verpt+1]
		//line ver.go.y:91
		{
			verVAL.tmp.Version = verDollar[3].val
			verVAL.tmp.Operation = "ST"
		}
	case 20:
		verDollar = verS[verpt-3 : verpt+1]
		//line ver.go.y:96
		{
			verVAL.tmp.Version = verDollar[3].val
			verVAL.tmp.Operation = "STE"
		}
	}
	goto verstack /* stack new state and value */
}
