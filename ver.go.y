%{
package dpkg
%}

%union {
    val string
    info DepInfo
}

%token PKGNAME VERSION ARCH_SPEC PROFILE ARCH_QUALIFIER
%start all

%%

all:            packages
                {
                    saveResult(verlex, $1.info);
                }
        ;

packages:       group
                {
                    $$.info = $1.info;
                }
        |       group ',' packages
                {
                    var tt = $3.info
                    $$.info = $1.info;
                    $$.info.And = &tt;
                }
        ;

group:          pkg
                {
                    $$.info = $1.info;
                }
        |       pkg '|' group
                {
                    var tt = $3.info;
                    $$.info = $1.info;
                    $$.info.Or = &tt;
                }

pkg:
                pkgname
                {
                    $$.info.Name = $1.val;
                }
        |       pkg VERSION
                {
                    $$.info.VerMin = $2.val;
                }
        |       pkg ARCH_SPEC
                {
                    $$.info.Archs = getArrayString($2.val, " ");
                }
        |       pkg PROFILE
                {
                    $$.info.Profiles = getArrayString($2.val, " ");
                }
        ;

pkgname:        PKGNAME
        |       PKGNAME ARCH_QUALIFIER
        ;
%%
