//line parser_error.y:2
package lexer

import __yyfmt__ "fmt"

//line parser_error.y:2
import "os"
import "fmt"

//line parser_error.y:7
type errorSymType struct {
	yys         int
	expecting   []string
	unexpecting string
	val         string
}

const EXPECTING = 57346
const HEAD = 57347
const ORinError = 57348
const COMMA = 57349
const TOKEN = 57350

var errorToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"EXPECTING",
	"HEAD",
	"ORinError",
	"COMMA",
	"TOKEN",
}
var errorStatenames = [...]string{}

const errorEofCode = 1
const errorErrCode = 2
const errorInitialStackSize = 16

//line parser_error.y:49

type myLexer struct {
	*SimpleLexer
	errorSymType
}

func (l myLexer) Lex(v *errorSymType) int {
	t, str := l.Token()
	v.val = str
	return int(t)
}

func forceSetmyLexerResult(l errorLexer, r errorSymType) {
	l.(*myLexer).errorSymType = r
}

func (l myLexer) Error(e string) {
	fmt.Fprintf(os.Stderr, "Internal error: %v when parsing %q\n", e, l.Input)
}

func ParseError(e string, tokens []string) (string, []string) {
	errorErrorVerbose = true
	l := NewSimpleLexer("error lexer")

	l.AddIgnores(" ")
	l.Add(HEAD, []string{"syntax error: unexpected"})
	l.Add(COMMA, []string{","})
	l.Add(ORinError, []string{"or"})
	l.Add(EXPECTING, []string{"expecting"})
	l.Add(TOKEN, tokens)

	l.SetInputString(e)

	my := &myLexer{l, errorSymType{}}
	errorParse(my)

	return my.unexpecting, my.expecting
}

//line yacctab:1
var errorExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const errorNprod = 5
const errorPrivate = 57344

var errorTokenNames []string
var errorStates []string

const errorLast = 9

var errorAct = [...]int{

	9, 7, 3, 4, 8, 2, 5, 6, 1,
}
var errorPact = [...]int{

	0, -1000, -6, -4, 2, -7, -2, -1000, -8, -1000,
}
var errorPgo = [...]int{

	0, 8, 7,
}
var errorR1 = [...]int{

	0, 1, 1, 2, 2,
}
var errorR2 = [...]int{

	0, 5, 2, 3, 1,
}
var errorChk = [...]int{

	-1000, -1, 5, 8, 7, 4, -2, 8, 6, 8,
}
var errorDef = [...]int{

	0, -2, 0, 2, 0, 0, 1, 4, 0, 3,
}
var errorTok1 = [...]int{

	1,
}
var errorTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8,
}
var errorTok3 = [...]int{
	0,
}

var errorErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	errorDebug        = 0
	errorErrorVerbose = false
)

type errorLexer interface {
	Lex(lval *errorSymType) int
	Error(s string)
}

type errorParser interface {
	Parse(errorLexer) int
	Lookahead() int
}

type errorParserImpl struct {
	lval  errorSymType
	stack [errorInitialStackSize]errorSymType
	char  int
}

func (p *errorParserImpl) Lookahead() int {
	return p.char
}

func errorNewParser() errorParser {
	return &errorParserImpl{}
}

const errorFlag = -1000

func errorTokname(c int) string {
	if c >= 1 && c-1 < len(errorToknames) {
		if errorToknames[c-1] != "" {
			return errorToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func errorStatname(s int) string {
	if s >= 0 && s < len(errorStatenames) {
		if errorStatenames[s] != "" {
			return errorStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func errorErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !errorErrorVerbose {
		return "syntax error"
	}

	for _, e := range errorErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + errorTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := errorPact[state]
	for tok := TOKSTART; tok-1 < len(errorToknames); tok++ {
		if n := base + tok; n >= 0 && n < errorLast && errorChk[errorAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if errorDef[state] == -2 {
		i := 0
		for errorExca[i] != -1 || errorExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; errorExca[i] >= 0; i += 2 {
			tok := errorExca[i]
			if tok < TOKSTART || errorExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if errorExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += errorTokname(tok)
	}
	return res
}

func errorlex1(lex errorLexer, lval *errorSymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = errorTok1[0]
		goto out
	}
	if char < len(errorTok1) {
		token = errorTok1[char]
		goto out
	}
	if char >= errorPrivate {
		if char < errorPrivate+len(errorTok2) {
			token = errorTok2[char-errorPrivate]
			goto out
		}
	}
	for i := 0; i < len(errorTok3); i += 2 {
		token = errorTok3[i+0]
		if token == char {
			token = errorTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = errorTok2[1] /* unknown char */
	}
	if errorDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", errorTokname(token), uint(char))
	}
	return char, token
}

func errorParse(errorlex errorLexer) int {
	return errorNewParser().Parse(errorlex)
}

func (errorrcvr *errorParserImpl) Parse(errorlex errorLexer) int {
	var errorn int
	var errorVAL errorSymType
	var errorDollar []errorSymType
	_ = errorDollar // silence set and not used
	errorS := errorrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	errorstate := 0
	errorrcvr.char = -1
	errortoken := -1 // errorrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		errorstate = -1
		errorrcvr.char = -1
		errortoken = -1
	}()
	errorp := -1
	goto errorstack

ret0:
	return 0

ret1:
	return 1

errorstack:
	/* put a state and value onto the stack */
	if errorDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", errorTokname(errortoken), errorStatname(errorstate))
	}

	errorp++
	if errorp >= len(errorS) {
		nyys := make([]errorSymType, len(errorS)*2)
		copy(nyys, errorS)
		errorS = nyys
	}
	errorS[errorp] = errorVAL
	errorS[errorp].yys = errorstate

errornewstate:
	errorn = errorPact[errorstate]
	if errorn <= errorFlag {
		goto errordefault /* simple state */
	}
	if errorrcvr.char < 0 {
		errorrcvr.char, errortoken = errorlex1(errorlex, &errorrcvr.lval)
	}
	errorn += errortoken
	if errorn < 0 || errorn >= errorLast {
		goto errordefault
	}
	errorn = errorAct[errorn]
	if errorChk[errorn] == errortoken { /* valid shift */
		errorrcvr.char = -1
		errortoken = -1
		errorVAL = errorrcvr.lval
		errorstate = errorn
		if Errflag > 0 {
			Errflag--
		}
		goto errorstack
	}

errordefault:
	/* default state action */
	errorn = errorDef[errorstate]
	if errorn == -2 {
		if errorrcvr.char < 0 {
			errorrcvr.char, errortoken = errorlex1(errorlex, &errorrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if errorExca[xi+0] == -1 && errorExca[xi+1] == errorstate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			errorn = errorExca[xi+0]
			if errorn < 0 || errorn == errortoken {
				break
			}
		}
		errorn = errorExca[xi+1]
		if errorn < 0 {
			goto ret0
		}
	}
	if errorn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			errorlex.Error(errorErrorMessage(errorstate, errortoken))
			Nerrs++
			if errorDebug >= 1 {
				__yyfmt__.Printf("%s", errorStatname(errorstate))
				__yyfmt__.Printf(" saw %s\n", errorTokname(errortoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for errorp >= 0 {
				errorn = errorPact[errorS[errorp].yys] + errorErrCode
				if errorn >= 0 && errorn < errorLast {
					errorstate = errorAct[errorn] /* simulate a shift of "error" */
					if errorChk[errorstate] == errorErrCode {
						goto errorstack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if errorDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", errorS[errorp].yys)
				}
				errorp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if errorDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", errorTokname(errortoken))
			}
			if errortoken == errorEofCode {
				goto ret1
			}
			errorrcvr.char = -1
			errortoken = -1
			goto errornewstate /* try again in the same state */
		}
	}

	/* reduction by production errorn */
	if errorDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", errorn, errorStatname(errorstate))
	}

	errornt := errorn
	errorpt := errorp
	_ = errorpt // guard against "declared and not used"

	errorp -= errorR2[errorn]
	// errorp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if errorp+1 >= len(errorS) {
		nyys := make([]errorSymType, len(errorS)*2)
		copy(nyys, errorS)
		errorS = nyys
	}
	errorVAL = errorS[errorp+1]

	/* consult goto table to find next state */
	errorn = errorR1[errorn]
	errorg := errorPgo[errorn]
	errorj := errorg + errorS[errorp].yys + 1

	if errorj >= errorLast {
		errorstate = errorAct[errorg]
	} else {
		errorstate = errorAct[errorj]
		if errorChk[errorstate] != -errorn {
			errorstate = errorAct[errorg]
		}
	}
	// dummy call; replaced with literal code
	switch errornt {

	case 1:
		errorDollar = errorS[errorpt-5 : errorpt+1]
		//line parser_error.y:26
		{
			errorVAL.unexpecting = errorDollar[2].val
			errorVAL.expecting = errorDollar[5].expecting
			forceSetmyLexerResult(errorlex, errorVAL)
		}
	case 2:
		errorDollar = errorS[errorpt-2 : errorpt+1]
		//line parser_error.y:32
		{
			errorVAL.unexpecting = errorDollar[2].val
			forceSetmyLexerResult(errorlex, errorVAL)
		}
	case 3:
		errorDollar = errorS[errorpt-3 : errorpt+1]
		//line parser_error.y:40
		{
			errorVAL.expecting = append(errorDollar[1].expecting, errorDollar[3].val)
		}
	case 4:
		errorDollar = errorS[errorpt-1 : errorpt+1]
		//line parser_error.y:44
		{
			errorVAL.expecting = append(errorVAL.expecting, errorDollar[1].val)
		}
	}
	goto errorstack /* stack new state and value */
}
