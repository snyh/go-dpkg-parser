package dpkg

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestListVersion(t *testing.T) {
	vs := []string{
		"python:native",
		"debhelper (>= 3.10), gcc (<< 7.0)",
		"debhelper (>= 3.10)",
		"gcc,wget",
		"gcc",
		"gcc|clang",
		"python-gdbm (>= 2.4.3)",
		"abc-dev [amd64]",
		"libbabeltrace-dev [amd64 armel armhf i386 kfreebsd-amd64 kfreebsd-i386 mips mipsel mips64el powerpc s390x]",
		"libghc-aeson-qq-dev [!mips !i386]",
	}
	for _, v := range vs {
		r, err := ParseDepends(v)
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("%q --> %+v", v, r)
	}
}
func TestAllDepends(t *testing.T) {
	bs, err := ioutil.ReadFile("testdata/deps.list")
	if err != nil {
		t.Fatal("Can't load deps.list", err)
	}

	for _, v := range strings.Split(string(bs), "\n") {
		_, err := ParseDepends(v)
		if err != nil {
			t.Fatal(err)
		}
	}
}
