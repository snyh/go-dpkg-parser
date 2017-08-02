%{
package dpkg
%}


%token ANY NATIVE

%token ','
%token '!'
%token ')' '('
%token ']' '['
%token '<' '>'
%token ALPHA_NUMERIC
%token PROFILE
%token ARCH_NAME

%left '.' '+' '~'

%union {
    val string
    r []Depend
    tmp Depend
}

%start expr

%%

expr:
                        group
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
