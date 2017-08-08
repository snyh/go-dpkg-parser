package dpkg

import (
	"bufio"
	"bytes"
	"fmt"
	"reflect"
	"runtime"
	"strings"
	"testing"
)

func TestParsingDBComponent(t *testing.T) {
	cs, err := LoadPackages("testdata/Packages")
	Assert(t, err, nil)
	Assert(t, len(cs), 523)
}

func TestParseLine(t *testing.T) {
	fs := getArrayString("haskell-mbox (0.3.3-3)", " ")
	Assert(t, len(fs), 2)
	Assert(t, fs[0], "haskell-mbox")
}

func TestBuildCache(t *testing.T) {
	targetDir := "testdata"
	suite := "unstable"

	m, err := NewSuite("http://packages.deepin.com/deepin/", suite, targetDir, "")
	Assert(t, err, nil)

	p, err := m.Archives["amd64"].FindBinary("lastore-daemon")
	Assert(t, err, nil)

	Assert(t, p.Homepage, "http://github.com/linuxdeepin/lastore-daemon")
}

func TestLarageControlFile(t *testing.T) {
	largeControlFile := fmt.Sprintf(testBinary, strings.Repeat("t", bufio.MaxScanTokenSize))

	b := bytes.NewBufferString(largeControlFile)
	ts, err := NewControlFiles(b, ScanBufferSize)
	Assert(t, err, nil)
	Assert(t, len(ts), 1)

	c, err := NewControlFile(largeControlFile)
	Assert(t, err, nil)
	p, err := c.ToBinary()
	Assert(t, err, nil)
	Assert(t, len(c.Get("Test")), bufio.MaxScanTokenSize)
	Assert(t, p.Package, "aac-enc")
}

func TestBinaryPackage(t *testing.T) {
	p := buildTestPackageBinary(t, testBinary)
	Assert(t, p.Filename, "pool/non-free/f/fdk-aac/aac-enc_0.1.3+20140816-2_amd64.deb")
	Assert(t, p.Size, 666554)

	Assert(t, p.Source, "fdk-aac")
	Assert(t, p.SourceVersion, "30")

	p = buildTestPackageBinary(t, testBinary2)
	Assert(t, p.Package, "apps.com.txsp")
	Assert(t, p.Source, p.Package)

	p = buildTestPackageBinary(t, testMultline)
	Assert(t, p.Package, "lib32gcc-4.9-dev")
}

func TestGetArrayString(t *testing.T) {
	Assert(t, getArrayString("a  b", " "), []string{"a", "b"})
}

func Assert(t *testing.T, left interface{}, right interface{}) {
	if !reflect.DeepEqual(left, right) {
		pc, file, line, ok := runtime.Caller(1)
		if ok {
			f := runtime.FuncForPC(pc)
			t.Fatalf("%v(%T) != %v(%T) \n\tat %s:%d:%s", left, left, right, right, file, line, f.Name())
		} else {
			t.Fatalf("%v(%T) != %v(%T)", left, left, right, right)
		}
	}
}

func buildTestPackageSource(t *testing.T, v string) SourcePackage {
	cf, err := NewControlFile(v)
	if err != nil {
		t.Fatal(err)
	}
	p, err := cf.ToSource()
	if err != nil {
		t.Fatal(err)
	}
	return p
}
func buildTestPackageBinary(t *testing.T, v string) BinaryPackage {
	cf, err := NewControlFile(v)
	if err != nil {
		t.Fatal(err)
	}
	p, err := cf.ToBinary()
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestDSC(t *testing.T) {
	p := buildTestPackageSource(t, testDSC)
	Assert(t, p.Binary, []string{"aften", "libaften0", "libaften-dev"})
	Assert(t, p.Format, "1.0")

	p = buildTestPackageSource(t, testNonePackList)
	Assert(t, p.PackageList[0], PackageListItem{Name: p.Package, Ptype: "deb", Archs: []string{"abc"}})
}

func TestControlFile(t *testing.T) {
	d, err := NewControlFile(testBinary)
	Assert(t, err, nil)

	Assert(t, d.Get("Package"), "aac-enc")

	Assert(t, d.Get("Source"), "fdk-aac   (30)")

	Assert(t, d.Get("Version"), "0.1.3+20140816-2")

	Assert(t, d.Get("installed-size"), "705")

	Assert(t, d.Get("archiTecTure"), "amd64")

	Assert(t, d.Get("depends"), "libfdk-aac0 (= 0.1.3+20140816-2), libc6 (>= 2.4)")

	Assert(t, d.Get("description"), `Fraunhofer FDK AAC Codec Library - frontend binary
 test multiline`)

	Assert(t, d.Get("priority"), "optional")

	Assert(t, d.Get("Filename"), "pool/non-free/f/fdk-aac/aac-enc_0.1.3+20140816-2_amd64.deb")
}

var testBinary = `
Package: aac-enc
Source: fdk-aac   (30)
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

var testMultline = `
Package: lib32gcc-4.9-dev
Source: gcc-4.9
Version: 4.9.4-2
Installed-Size: 6211
Maintainer: Debian GCC Maintainers <debian-gcc@lists.debian.org>
Architecture: amd64
Depends: gcc-4.9-base (= 4.9.4-2), lib32gcc1 (>= 1:4.9.4-2), libx32gcc1 (>= 1:4.9.4-2), lib32gomp1 (>= 4.9.4-2), libx32gomp1 (>= 4.9.4-2\
), lib32itm1 (>= 4.9.4-2), libx32itm1 (>= 4.9.4-2), lib32atomic1 (>= 4.9.4-2), libx32atomic1 (>= 4.9.4-2), lib32asan1 (>= 4.9.4-2), libx\
32asan1 (>= 4.9.4-2), lib32ubsan0 (>= 4.9.4-2), libx32ubsan0 (>= 4.9.4-2), lib32cilkrts5 (>= 4.9.4-2), libx32cilkrts5 (>= 4.9.4-2), lib3\
2quadmath0 (>= 4.9.4-2), libx32quadmath0 (>= 4.9.4-2)
Recommends: libc6-dev (>= 2.13-5)
Size: 1903246
SHA256: fa82adbc224d01541594bd1b11996288cf5fe73b930d4850d795a75bda4d170a
SHA1: 019b774f9db5db7812dac08189cd744b84a01631
MD5sum: bb3f9843fa97cc5fdc30ca934cada406
Description: GCC support library (32 bit development files)
 This package contains the headers and static library files necessary for
 building C programs which use libgcc, libgomp, libquadmath, libssp or libitm.
Homepage: http://gcc.gnu.org/
Tag: devel::library, role::devel-lib
Section: libdevel
Priority: optional
Filename: pool/main/g/gcc-4.9/lib32gcc-4.9-dev_4.9.4-2_amd64.deb
`

const testNonePackList = `
Package: tmpreaper
Binary: tmpreaper
Version: 1.6.13+nmu1
Maintainer: Paul Slootman <paul@debian.org>
Build-Depends: debhelper (>= 5), e2fslibs-dev, po-debconf
Architecture: abc
Standards-Version: 3.8.3.0
Format: 1.0
Files:
 2e628b122fae3896cafb9ea31887021c 748 tmpreaper_1.6.13+nmu1.dsc
 36bffb38fbdd28b9de8af229faabf5fe 141080 tmpreaper_1.6.13+nmu1.tar.gz
Checksums-Sha1:
 ea9c60662bb8998e9486204c0a41e7bb003155c4 748 tmpreaper_1.6.13+nmu1.dsc
 96a490a9c2df6d3726af8df299e5aedd7d49fbfe 141080 tmpreaper_1.6.13+nmu1.tar.gz
Checksums-Sha256:
 8782f6fcdf98ba2f77ee278d13806d3e4f7e0c991a8940473fa0afe1b7e466f9 748 tmpreaper_1.6.13+nmu1.dsc
 c88f05b5d995b9544edb7aaf36ac5ce55c6fac2a4c21444e5dba655ad310b738 141080 tmpreaper_1.6.13+nmu1.tar.gz
Directory: pool/main/t/tmpreaper
Priority: source
Section: admin
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

const testBinary2 = `
Package: apps.com.txsp
Version: 1.1
Architecture: amd64
Maintainer: Deepin Packages Builder <packages@deepin.com>
Installed-Size: 23347
Depends: deepin-chrome-arc
Homepage: http://www.deepin.org
Priority: optional
Section: utils
Filename: pool/main/a/apps.com.txsp/apps.com.txsp_1.1_amd64.deb
Size: 22282694
SHA256: 91be41232545b999bd1a5805618f19b1853a333df7a3db4cee797789f9aa49c2
SHA1: abdffc72f4b91fba2b791c8f87a9e2b7cbae4a08
MD5sum: a50d7c5f974586f25ae7474aa103a571
Description: Tengxun Video for Chrome ARC
Description-md5: 1cc503533ab128384dee9c0dc7f5a689
`
