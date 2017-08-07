%{
package dpkg
%}

%union {
    val string
    info DepInfo
}

%token PKGNAME VERSION ARCH_QUALIFIER PROFILE
%token '|'

%start all

%%

all:            group
                {
                    saveResult(verlex, $1.info);
                }
        ;

group:          pkg
                {
                    $$.info = $1.info
                }
        |       pkg '|' group
                {
                    $$.info = $1.info;
                    $$.info.Or = &($3.info);
                }

pkg:
                PKGNAME
                {
                    $$.info.Name = $1.val;
                }
        |       pkg VERSION
                {
                    $$.info.VerMin = $2.val;
                }
        |       pkg ARCH_QUALIFIER
                {
                    $$.info.Archs = getArrayString($2.val, " ");
                }
        |       pkg PROFILE
                {
                    $$.info.Profiles = getArrayString($2.val, " ");
                }
        ;
%%
