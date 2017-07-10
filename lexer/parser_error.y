%{
package lexer
import "os"
import "fmt"
%}

%union {
    expecting []string
    unexpecting string
    val string
}

%token EXPECTING
%token HEAD
%token ORinError
%token COMMA

%token TOKEN

%start serror

%%

serror:
                HEAD TOKEN COMMA EXPECTING utokens
                {
                        $$.unexpecting = $2.val
                        $$.expecting = $5.expecting
                        forceSetmyLexerResult(errorlex, $$)
                }
        |       HEAD TOKEN
                {
                    $$.unexpecting = $2.val
                    forceSetmyLexerResult(errorlex, $$)
                }
                ;

utokens:
                utokens ORinError TOKEN
                {
                    $$.expecting = append($1.expecting, $3.val)
                }
        |       TOKEN
                {
                    $$.expecting = append($$.expecting, $1.val)
                }
                ;

%%

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

        my := &myLexer{l,errorSymType{}}
        errorParse(my)

        return my.unexpecting, my.expecting
}
