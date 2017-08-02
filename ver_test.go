package dpkg

import (
	"io/ioutil"
	"strings"
	"testing"
)

func TestParseDepends(t *testing.T) {
	vs := map[string]string{
		"valac":                                                   "valac",
		"libgtk-3-dev":                                            "libgtk-3-dev",
		"libgee-0.8-dev":                                          "libgee-0.8-dev",
		"libvte-2.91-dev":                                         "libvte-2.91-dev",
		"libjson-glib-dev":                                        "libjson-glib-dev",
		"libsecret-1-dev":                                         "libsecret-1-dev",
		"libwnck-3-dev":                                           "libwnck-3-dev",
		"cmake":                                                   "cmake",
		"debhelper (>= 3.10)":                                     "debhelper",
		"libbabeltrace-dev [amd64 mipsel mips64el powerpc s390x]": "libbabeltrace-dev",
		"libghc-aeson-qq-dev [!mips !i386]":                       "libghc-aeson-qq-dev",
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
		_, err := parseDepends("dump (>=" + v + ")")
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
		_, err := parseDepends(v)
		if err != nil {
			t.Fatal(err)
		}
	}
}
