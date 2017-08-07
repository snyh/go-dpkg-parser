package dpkg

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestDepMatch(t *testing.T) {
	str := "libc6.1-dev   (>= 2.13-5)    [  alpha ia64] | ulibc"
	info, err := parseDepInfo(str)
	Assert(t, err, nil)
	Assert(t, info.Name, "libc6.1-dev")
	Assert(t, info.Archs, []string{"alpha", "ia64"})
	Assert(t, info.Match("amd64", ""), false)
	Assert(t, info.Or.Name, "ulibc")

	str = "libbabeltrace-dev [amd64]"
	info, err = parseDepInfo(str)
	Assert(t, err, nil)
	Assert(t, info.Archs, []string{"amd64"})

}
func TestParseDepends(t *testing.T) {
	vs := map[string]string{
		"valac":                                                          "valac",
		"libgtk-3-dev":                                                   "libgtk-3-dev",
		"libgee-0.8-dev":                                                 "libgee-0.8-dev",
		"libvte-2.91-dev":                                                "libvte-2.91-dev",
		"libjson-glib-dev":                                               "libjson-glib-dev",
		"libsecret-1-dev":                                                "libsecret-1-dev",
		"libwnck-3-dev":                                                  "libwnck-3-dev",
		"cmake":                                                          "cmake",
		"debhelper (>= 3.10)":                                            "debhelper",
		"libbabeltrace-dev [amd64 mipsel mips64el powerpc s390x]":        "libbabeltrace-dev",
		"libghc-aeson-qq-dev [!mips !i386]":                              "libghc-aeson-qq-dev",
		"g++-multilib [amd64 i386 powerpc ppc64 s390 sparc] <!nobiarch>": "g++-multilib",
	}

	for str, rightName := range vs {
		info, err := parseDepInfo(str)
		Assert(t, err, nil)
		Assert(t, info.Name, rightName)
	}
}

func TestAllVersions(t *testing.T) {
	t.Skip()
	bs, err := ioutil.ReadFile("testdata/ver.list")
	if err != nil {
		t.Fatal("Can't load deps.list", err)
	}

	for _, v := range strings.Split(string(bs), "\n") {
		_, err := parseDepInfo("dump (>=" + v + ")")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAllDepends(t *testing.T) {
	t.Skip()
	bs, err := ioutil.ReadFile("testdata/deps.list")
	if err != nil {
		t.Fatal("Can't load deps.list", err)
	}

	for _, v := range strings.Split(string(bs), "\n") {
		_, err := parseDepInfo(v)
		if err != nil {
			t.Fatal(err)
		}
	}
}
