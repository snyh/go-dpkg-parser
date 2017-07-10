%{
package dpkg
import "./lexer"
%}


%token ANY NATIVE

%token ','
%token '!'
%token ')' '('
%token ']' '['
%token '<' '>'
%token VER_NUM
%token PKG_NAME
%token PROFILE
%token ARCH_NAME

%union {
    val string
    r []Depend
    tmp Depend
}

%start groups

%%

groups:
                group ',' groups
                {
                    $$.r = append([]Depend{$1.tmp}, $3.r...)
                    saveResult(verlex, $$.r)
                }
        |       group
                {
                    $$.r = []Depend{$1.tmp}
                    saveResult(verlex, $$.r)
                }
        ;

group:
                package '|' group
                {
                    $$.r = append([]Depend{$1.tmp}, $3.r...)
                }
        |       package
                {
                    $$.r = []Depend{$1.tmp}
                }
        ;

package:        PKG_NAME
                {
                    $$.tmp.Name = $1.val
                    $$.tmp.Version = ""
                    $$.tmp.Operation = ""
                }
        |       package ':' arch_qualifier
        |       package '(' ver_spec ')'
                {
                    $$.tmp.Name = $1.val
                    $$.tmp.Version = $3.tmp.Version
                    $$.tmp.Operation = $3.tmp.Operation
                }
        |       package '[' arch_spec ']'
        |       package '<' restriction_formula '>'
        ;

arch_qualifier: ANY
        |       NATIVE
        |       ARCH_NAME
        ;

ver_spec:       VER_NUM
                {
                    $$.tmp.Version = $1.val
                }
        |       '>' '>' VER_NUM
                {
                    $$.tmp.Version = $3.val
                    $$.tmp.Operation = "GT"
                }
        |       '>' '=' VER_NUM
               {
                    $$.tmp.Version = $3.val
                    $$.tmp.Operation = "GTE"
                }
        |       '<' '<' VER_NUM
                {
                    $$.tmp.Version = $3.val
                    $$.tmp.Operation = "ST"
                }
        |       '<' '=' VER_NUM
                {
                    $$.tmp.Version = $3.val
                    $$.tmp.Operation = "STE"
                }
        ;
arch_spec:      ARCH_NAME
        |       '!' arch_spec
        |       ARCH_NAME arch_spec
        ;

restriction_formula:
                PROFILE
        |       '!' PROFILE
        ;

%%

type myLexer struct {
    e string
    slex *lexer.SimpleLexer
    r []Depend
}

func NewMyLexer(raw string, archList, profileList []string) *myLexer {
    sl := lexer.NewSimpleLexer("dsc-version")
    sl.SetInputString(raw)

    sl.AddIgnores(" ")

    sl.Add(ANY, []string{"any"})
    sl.Add(NATIVE, []string{"native"})

    sl.SetBasicToken("()[]<>=!,")

    sl.AddR(VER_NUM, []string{})
    sl.AddI(ARCH_NAME, archList);
    sl.AddI(PROFILE, profileList);
    sl.AddR(PKG_NAME, []string{`^[a-z][a-z0-9\.\-\+]+`})
    sl.AddR(VER_NUM, []string{`^[0-9]+[a-z0-9\.\+\-\~\:]*`})
    return &myLexer{
       slex: sl,
    }
}

func saveResult(l verLexer, r []Depend) {
    l.(*myLexer).r =  r
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
    Name string
    Version string
    Operation string
}

func ParseDepends(v string) ([]Depend, error) {
    archList := []string{"amd64", "i386", "ia64",
                         "mips", "mipsel", "mips64el","mips64", "mipsn32", "mipsn32el",
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
