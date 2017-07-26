package dpkg

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParsingDBComponent(t *testing.T) {
	cs, err := LoadControlFileGroup("testdata/Packages")
	Assert(t, err, nil)
	Assert(t, len(cs), 523)
}

func TestBuildCache(t *testing.T) {
	targetDir := "testdata"
	codeName := "unstable"

	m, err := NewSuite("http://pools.corp.deepin.com/deepin/", codeName, targetDir)
	Assert(t, err, nil)

	p, err := m.FindBinary("lastore-daemon", "all")
	Assert(t, err, nil)

	Assert(t, p.Homepage, "http://github.com/linuxdeepin/lastore-daemon")
}

func TestLarageControlFile(t *testing.T) {
	largeControlFile := fmt.Sprintf(testBinary, strings.Repeat("t", bufio.MaxScanTokenSize))
	b := bytes.NewBufferString(largeControlFile)
	ts, err := ParseControlFileGroup(b)
	if err != nil {
		t.Fatal(err)
	}
	if len(ts) != 1 {
		t.Fatal("Parsing failed, token size is too long")
	}
}

func TestBinaryPackage(t *testing.T) {
	cf, err := NewControlFile(bytes.NewBuffer([]byte(testBinary)))
	if err != nil {
		t.Fatal(err)
	}
	p, err := cf.ToBinary()
	if err != nil {
		t.Fatal(err)
	}
	Assert(t, p.Filename, "pool/non-free/f/fdk-aac/aac-enc_0.1.3+20140816-2_amd64.deb")
	Assert(t, p.Size, 666554)
}

func Assert(t *testing.T, left interface{}, right interface{}) {
	if !reflect.DeepEqual(left, right) {
		t.Fatalf("%v(%T) != %v(%T)", left, left, right, right)
	}
}

func TestDSC(t *testing.T) {
	cf, err := NewControlFile(bytes.NewBuffer([]byte(testDSC)))
	if err != nil {
		t.Fatal(err)
	}
	p, err := cf.ToSource()
	if err != nil {
		t.Fatal(err)
	}
	Assert(t, p.Binary, []string{"aften", "libaften0", "libaften-dev"})
	Assert(t, p.Format, "1.0")
}

func TestControlFile(t *testing.T) {
	d, err := NewControlFile(bytes.NewBuffer([]byte(testBinary)))
	Assert(t, err, nil)

	Assert(t, d.GetString("Package"), "aac-enc")

	Assert(t, d.GetString("Source"), "fdk-aac")

	Assert(t, d.GetString("Version"), "0.1.3+20140816-2")

	Assert(t, d.GetString("installed-size"), "705")

	Assert(t, d.GetString("archiTecTure"), "amd64")

	Assert(t, d.GetString("depends"), "libfdk-aac0 (= 0.1.3+20140816-2), libc6 (>= 2.4)")

	Assert(t, d.GetString("description"), `Fraunhofer FDK AAC Codec Library - frontend binary
 test multiline`)

	Assert(t, d.GetString("priority"), "optional")

	Assert(t, d.GetString("Filename"), "pool/non-free/f/fdk-aac/aac-enc_0.1.3+20140816-2_amd64.deb")
}

var testBinary = `
Package: aac-enc
Source: fdk-aac
Version: 0.1.3+20140816-2
Installed-Size: 705
Maintainer: Debian Multimedia Maintainers <pkg-multimedia-maintainers@lists.alioth.debian.org>
Architecture: amd64
Depends: libfdk-aac0 (= 0.1.3+20140816-2), libc6 (>= 2.4)
Size: 666554
SHA256: d09f8c35f8817bc67b67ebc7af94d7b26ba656af2bea4ed579e13e03db718cee
SHA1: b9a70c3b65f7ad6b62f56c2b8cc916b156c38713
MD5sum: 9703f7d0d4463b198bfd57b45fefd8ab
Description: Fraunhofer FDK AAC Codec Library - frontend binary
 test multiline
Homepage: https://github.com/mstorsjo/fdk-aac
Description-md5: 16f812d0c8b3e09448f6f7d88536e135
Section: non-free/sound
Priority: optional
Filename: pool/non-free/f/fdk-aac/aac-enc_0.1.3+20140816-2_amd64.deb
Test: %s
`

var testDSC = `Package: aften
Format: 1.0
Binary: aften, libaften0, libaften-dev
Architecture: any
Version: 0.0.8svn20100103-0.0
Maintainer: Christian Marillat <marillat@debian.org>
Homepage: http://aften.sourceforge.net/
Standards-Version: 3.8.3
Build-Depends: debhelper (>= 7), cmake
Package-List:
 aften deb sound extra arch=any
 libaften-dev deb libdevel extra arch=any
 libaften0 deb libs extra arch=any
Priority: extra
Section: sound
Directory: pool/contrib/a/aften
Files:
 2ff61aad7ea2818cc4cff61e5f10310b 964 aften_0.0.8svn20100103-0.0.dsc
 14f13a6a3b9489ac6542105fb7c3eed6 125272 aften_0.0.8svn20100103.orig.tar.gz
 de015da7faef5886ba6b8fd2ec09709b 2313 aften_0.0.8svn20100103-0.0.diff.gz
Checksums-Sha1:
 d0f80e1533092768f086e48c34286cb5feb61274 964 aften_0.0.8svn20100103-0.0.dsc
 ef8980898599a555dfd5b0565f6033ab516ac0ec 125272 aften_0.0.8svn20100103.orig.tar.gz
 d44082ce23d8907ee020148b3c27fb7d3b6ddf1f 2313 aften_0.0.8svn20100103-0.0.diff.gz
Checksums-Sha256:
 204b23580d54928073c76453ace942ee0731d90214062e48236ac2ebf1ba46c4 964 aften_0.0.8svn20100103-0.0.dsc
 372643b7b62258504f80c73d75a59e3cbfc5dec6f9e8bda626a9b5816f269c8f 125272 aften_0.0.8svn20100103.orig.tar.gz
 2e78c5eb9dd74cd1f6e08a76fdd3c1633fd8f43ecfafe36783b0b0ed47f784f6 2313 aften_0.0.8svn20100103-0.0.diff.gz
`
