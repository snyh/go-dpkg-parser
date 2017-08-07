package dpkg

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestFilterMatch(t *testing.T) {
	str := "libbabeltrace-dev [amd64], libabc [m] | lib0 [mips] |lib2,   libefg [xx], lib3"
	info, err := ParseDepInfo(str)
	Assert(t, err, nil)
	info, err = info.Filter("i386", "")
	Assert(t, err, nil)
	Assert(t, info.String(), "lib2, lib3")
}
func TestDepFilter(t *testing.T) {
	str := "libbabeltrace-dev [amd64], libabc [m] | lib0 [mips] |lib2,   libefg, lib3"
	info, err := ParseDepInfo(str)
	Assert(t, err, nil)

	info, err = info.Filter("i386", "")
	Assert(t, err, nil)
	Assert(t, info.String(), "lib2, libefg, lib3")
}

func TestDepMatch(t *testing.T) {
	str := `binutils:native (>= 2.27.90.20170221) | binutils-multiarch:native (>= 2.27.90.20170221)`
	info, err := ParseDepInfo(str)
	Assert(t, info.Name, "binutils")
	Assert(t, info.Or != nil, true)

	str = "libc6.1-dev   (>= 2.13-5)    [  alpha ia64] | ulibc | linuxc"
	info, err = ParseDepInfo(str)
	Assert(t, err, nil)
	Assert(t, info.Name, "libc6.1-dev")
	Assert(t, info.Archs, []string{"alpha", "ia64"})
	Assert(t, info.match("amd64", ""), false)
	Assert(t, info.Or.Name, "ulibc")
	Assert(t, info.Or.Or.Name, "linuxc")

	str = "libbabeltrace-dev [amd64], libabc | lib0 |lib2,   libefg, lib3"
	info, err = ParseDepInfo(str)

	Assert(t, err, nil)
	Assert(t, info.Archs, []string{"amd64"})

	Assert(t, info.And.Name, "libabc")
	Assert(t, info.And.Or.Name, "lib0")
	Assert(t, info.And.Or.Or.Name, "lib2")

	Assert(t, info.And.And.Name, "libefg")
	Assert(t, info.And.And.And.Name, "lib3")

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
		info, err := ParseDepInfo(str)
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
		_, err := ParseDepInfo("dump (>=" + v + ")")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func TestAllDepends(t *testing.T) {
	bs, err := ioutil.ReadFile("testdata/deps.list")
	if err != nil {
		t.Fatal("Can't load deps.list", err)
	}

	for _, v := range strings.Split(string(bs), "\n") {
		info, err := ParseDepInfo(v)
		if err != nil {
			t.Fatal(err)
		}
		t.Log(v, "-->", info)
		c1, c2 := strings.Count(v, ","), strings.Count(v, "|")
		cc1, cc2 := countAndAndOr(info)
		Assert(t, c1, cc1)
		Assert(t, c2, cc2)
	}
}

func countAndAndOr(info DepInfo) (int, int) {
	and, or := 0, 0
	if info.And != nil {
		and++
		t1, t2 := countAndAndOr(*info.And)
		and += t1
		or += t2
	}
	if info.Or != nil {
		or++
		t1, t2 := countAndAndOr(*info.Or)
		and += t1
		or += t2
	}
	return and, or
}
