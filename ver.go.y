%{
package dpkg
%}

%union {
    val string
    info DepInfo
}

%token PKGNAME VERSION ARCH_QUALIFIER PROFILE

%start all

%%

all:  pkg
                {
                    saveResult(verlex, $1.info);
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
