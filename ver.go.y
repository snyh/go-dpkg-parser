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
%token ALPHA_NUMERIC
%token PKG_NAME
%token PROFILE
%token ARCH_NAME

%left ',' '.' '+' '~'
%right '-'

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

pkg_name:       ALPHA_NUMERIC
        |       pkg_name '.' ALPHA_NUMERIC
        |       pkg_name '-' ALPHA_NUMERIC
        |       pkg_name '+' ALPHA_NUMERIC
        ;

package:        pkg_name
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

ver_spec:       ver_num
                {
                    $$.tmp.Version = $1.val
                }
        |       '>' '>' ver_num
                {
                    $$.tmp.Version = $3.val
                    $$.tmp.Operation = "GT"
                }
        |       '>' '=' ver_num
               {
                    $$.tmp.Version = $3.val
                    $$.tmp.Operation = "GTE"
                }
        |       '<' '<' ver_num
                {
                    $$.tmp.Version = $3.val
                    $$.tmp.Operation = "ST"
                }
        |       '<' '=' ver_num
                {
                    $$.tmp.Version = $3.val
                    $$.tmp.Operation = "STE"
                }
        ;

ver_num:        upstream_version
        |       upstream_version2 '-' debian_version
        ;

upstream_version:
                ALPHA_NUMERIC
        |       upstream_version '+' upstream_version
        |       upstream_version '.' upstream_version
        |       upstream_version '~' upstream_version
        ;

upstream_version2:
                upstream_version
        |       upstream_version2 '-' upstream_version2
        ;

debian_version:
                ALPHA_NUMERIC
        |       ALPHA_NUMERIC '+' debian_version
        |       ALPHA_NUMERIC '.' debian_version
        |       ALPHA_NUMERIC '~' debian_version
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

    sl.SetBasicToken("()[]<>=!,:+-~.")

    sl.AddI(ARCH_NAME, archList);
    sl.AddI(PROFILE, profileList);
    sl.AddR(ALPHA_NUMERIC, []string{`^[0-9a-zA-Z]+`})
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
