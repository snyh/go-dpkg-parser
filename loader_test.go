package dpkg

import (
	"flag"
	"path"
	"strings"
	"testing"
)

var network = flag.Bool("network", false, "download test data from network")

func TestDumpRepository(t *testing.T) {
	repoURL := "http://pools.corp.deepin.com/deepin"
	rootDir := "/tmp/dump_repository"
	codeName := "unstable"

	cf, err := DownloadReleaseFile(repoURL, codeName)
	Assert(t, err, nil)

	rf, err := cf.ToReleaseFile()
	Assert(t, err, nil)

	_, err = DownloadRepository(repoURL, rf, path.Join(rootDir, rf.Hash))
	Assert(t, err, nil)
}

func TestHash(t *testing.T) {
	v := HashBytes([]byte("hello"))
	Assert(t, v, "5d41402abc4b2a76b9719d911017c592")
}

func TestRelease(t *testing.T) {
	cf, err := NewControlFile(testRelease)

	rf, err := cf.ToReleaseFile()
	Assert(t, err, nil)

	Assert(t, rf.CodeName, "experimental")

	Assert(t, len(rf.Architectures), 1)
	Assert(t, rf.Architectures[0], "amd64")
	Assert(t, strings.Join(rf.Components, ""), "non-free")

	Assert(t, len(rf.fileInfos), 31)
	Assert(t, len(rf.FileInfos()), 2)
	pf := rf.fileInfos[2]
	Assert(t, pf.Size, uint64(0x8f))
	Assert(t, pf.MD5, "f23e539f4e40f8491b5b5512d1e7aaa9")
	Assert(t, pf.Path, "main/binary-i386/Release")
}

const testRelease = `
Origin: deepin
Label: deepin
Codename: experimental
Version: 151218
Date: Tue, 02 Feb 2016 09:14:25 UTC
Architectures: amd64
Components: non-free
Description: deepin 2015 packages
 deepin 2015 mirrors for test
MD5Sum:
 5b8493e6afa5ccd3dd7bf809c9eb7dd7 49843895 main/binary-i386/Packages
 246d846a8137550a65645fa65ca08e25 13155202 main/binary-i386/Packages.gz
 f23e539f4e40f8491b5b5512d1e7aaa9 143 main/binary-i386/Release
 3566ea908e026a5bcff97bf56ddc00c6 531154 main/debian-installer/binary-i386/Packages
 0723d236434e392d1ea4c1d66672c827 128199 main/debian-installer/binary-i386/Packages.gz
 9a848962692eb8896df6242ff543e6a2 49972458 main/binary-amd64/Packages
 7c82a88d6887e1f85b517672a93ef80d 13170247 main/binary-amd64/Packages.gz
 71bee5809e10ee72fb96ff9267b073c9 144 main/binary-amd64/Release
 3c78b18b27daf5d7401cdf422b44e3b0 380273 main/debian-installer/binary-amd64/Packages
 ca9c2319a1dd542441487e9d0f9f370d 100654 main/debian-installer/binary-amd64/Packages.gz
 af5ba11d34ee3866e9ef26ed87d87e7b 37914073 main/source/Sources
 15e6bb3a0176a2882f7da899153bf44a 10497167 main/source/Sources.gz
 ae9bc1d9a7d7e2e79400c196d0553c3e 145 main/source/Release
 dba86955c7fe2878eb95e89a53aa558a 416317 contrib/binary-i386/Packages
 21b748c38c01e9c1d6673699ee2d9e3a 122174 contrib/binary-i386/Packages.gz
 83b77160c60535ee5fed32245aa6488e 146 contrib/binary-i386/Release
 bf13a82e9543b18c26e94dc4beeef580 416840 contrib/binary-amd64/Packages
 acd9c4355eba82421f3160a9be8d5c27 122304 contrib/binary-amd64/Packages.gz
 8d1a30565ff978554ed11e39d80c6ce2 147 contrib/binary-amd64/Release
 0548bb37c26802915cc98d838f1e45bb 278018 contrib/source/Sources
 25004c95d318e11c2c214f47bce1e692 84833 contrib/source/Sources.gz
 dd3b7e1bcbf2f9397d7be211056d43f3 148 contrib/source/Release
 0a491f77d55e486f51e309df6c4c456e 709369 non-free/binary-i386/Packages
 4c2e82a1938d5c2b787f57119b0a4aa7 175701 non-free/binary-i386/Packages.gz
 79fcb3e6a414f71d3db15215f50fe5e4 147 non-free/binary-i386/Release
 efd06a9bc2425c727626b249de1dd168 719937 non-free/binary-amd64/Packages
 9955efb933b72695cc927732fb273b43 179794 non-free/binary-amd64/Packages.gz
 b8f18141f566829e69d6daf6943c7732 148 non-free/binary-amd64/Release
 1c990f6a26ab7bdb6d98793daa9b3656 468876 non-free/source/Sources
 e8fa2b9c7ae4a0bc2efc6562e216b8c0 137694 non-free/source/Sources.gz
 9db073ef01ad70ee83a8cd522a9fc934 149 non-free/source/Release
SHA1:
 dafc1d4e3f1ae81c74476c08fafa6b9bf390be63 49843895 main/binary-i386/Packages
 7b6a592f57e20517248694bb1216753d6fee8537 13155202 main/binary-i386/Packages.gz
 d4cd2ac98dd8fe2efcfb0502a70e3ced753cfa2c 143 main/binary-i386/Release
 6d9cd14f4c480ede5a81f282ad3f2a0b0063a4f9 531154 main/debian-installer/binary-i386/Packages
 aa59041b274eed71ca88f40e53b95cd04c8bf634 128199 main/debian-installer/binary-i386/Packages.gz
 00e724ea84eec886c09cb656631b6ee8ff59fb6f 49972458 main/binary-amd64/Packages
 682f74589a61eed5919578ceac858cc4beaf551b 13170247 main/binary-amd64/Packages.gz
 2179b4e2c8b4f297e45941033fe3bb64a1e103ca 144 main/binary-amd64/Release
 3bb0f93da5c5db326be2b7bdb09199c5eacefedd 380273 main/debian-installer/binary-amd64/Packages
 12144e0ef609ae963cb0c8272945b648b6fe65da 100654 main/debian-installer/binary-amd64/Packages.gz
 0968f424bba7adbd87f01809943da30bee01571f 37914073 main/source/Sources
 0aeba47c5786bd9979f68f970342a382d9ff0d94 10497167 main/source/Sources.gz
 4fa719c63e6b2147cb147cce44024705ea4f955e 145 main/source/Release
 f3d250040a1ca9f282243112847e15d33749f5fe 416317 contrib/binary-i386/Packages
 b35e6fb3debaca83d09fbf1fa68ba3fb76b50b5f 122174 contrib/binary-i386/Packages.gz
 2ea1bfc091c9e03dd36d7f1a1595baebe89a92a1 146 contrib/binary-i386/Release
 375fc94f9400e9c80c1e09143739b3129eb6770d 416840 contrib/binary-amd64/Packages
 822a1410fc24c942a9f6a8c56d661da36cd14113 122304 contrib/binary-amd64/Packages.gz
 54b612275dcb186e5d88e43744ac3cbefe9bfea6 147 contrib/binary-amd64/Release
 aadb17fbdc4abbfbd3bf3134f52e0e603d483f12 278018 contrib/source/Sources
 eb65e9ba4638fc12734d41d00c9d637e3e3611b6 84833 contrib/source/Sources.gz
 d0d329eaec911bee8d12ed3f8bf78c2d6696dd5a 148 contrib/source/Release
 6221b12a49869b2fe27b007847eee59419a75f71 709369 non-free/binary-i386/Packages
 95d85b72d6f1c002e6ff6ceafaddd944c0580411 175701 non-free/binary-i386/Packages.gz
 47587e4845fd97977fbe6989ab083a5f0cdecf61 147 non-free/binary-i386/Release
 59c919d102d83330c18e229e039ecfb264c92a0e 719937 non-free/binary-amd64/Packages
 3021c3d9776eb6b3d737612d23d3033385855ae1 179794 non-free/binary-amd64/Packages.gz
 a99ade3bfd2712810edc9b054bf22f6c7faccc81 148 non-free/binary-amd64/Release
 91880149906f342bc3309a9aed2be9665931d684 468876 non-free/source/Sources
 02d6eb5bce498204561ce1e372978ff60cae78f7 137694 non-free/source/Sources.gz
 6f915b34dbe97143b33869f1bfcf26f34cd9b1ef 149 non-free/source/Release
SHA256:
 68780203b7ce2411ec1b57fc86d7353f7288128e1c4b870809088058e5ebd4ee 49843895 main/binary-i386/Packages
 ccb649869e510b7db1f7f0e945aabc9970f9aeb52efa5a5610dc7886d29232c3 13155202 main/binary-i386/Packages.gz
 9cbbd0e27b6665c4e7286f233d23fe7d2513508df810c5073fe1848c804f14a6 143 main/binary-i386/Release
 1181ccc37a71a156654814821f9f878f1de4133a21d9e047c562f75808ea089a 531154 main/debian-installer/binary-i386/Packages
 135f8a87c0bea4e299aa73297825e28b3324afc749913be42a9c7f8e196985f8 128199 main/debian-installer/binary-i386/Packages.gz
 43c7effe91f0ed5ced2cda93f2ea2cf4b567104daf26b1d8849acbb123ec3a69 49972458 main/binary-amd64/Packages
 066857398da76555bc5644658ba5806fee272114ee2596397ea053021ac357d7 13170247 main/binary-amd64/Packages.gz
 bbeaad42daf35a55644563e1577114ec562b7ebf7b93ede4a75ddc25602cb669 144 main/binary-amd64/Release
 26713c82a6ae61b503535295b51c55c0b0e48b511bcd96134758dcbbb4caa40c 380273 main/debian-installer/binary-amd64/Packages
 491b68c9ce558d5bb8484ae343ae41021e597367056acf831f13d9ef4ff8135e 100654 main/debian-installer/binary-amd64/Packages.gz
 1871be27164fc5e7638ea837ca1b9d7eda3c940e1c168fee8873832779b98c37 37914073 main/source/Sources
 3bef6735d96fba7808f82f7319fdc72613c823f1667a6b49e2be34023bb48ebf 10497167 main/source/Sources.gz
 ae3f8d1a5aec41ddc8fd1df92aa1ceb49ef709936972b82634ff30b8dbe3a48d 145 main/source/Release
 ef06efdc8169e8a849829c50cb4f00b83332c9f350d846c796db93677b292178 416317 contrib/binary-i386/Packages
 1e6e1f4d9a9ec10a2fa694b12bc437b626f680ca9396c73de6d079fbbc1b4c3a 122174 contrib/binary-i386/Packages.gz
 fee3c3a159d12ba17237704d4cc5a35d6af72002079816c355f02673e69fa430 146 contrib/binary-i386/Release
 1078fcf0c0a228e147c7d5f4139fa152dbce5f98767e5044bcb65d764c6d8e8d 416840 contrib/binary-amd64/Packages
 6c99ed35a819f482eaf8a864d797c933ec688aeabf3c28e74ad2dd922024b5fe 122304 contrib/binary-amd64/Packages.gz
 7132fa47a52bd79b6f793fa6c0c9c2da55010a4f2a8dee3d9a64d186872725de 147 contrib/binary-amd64/Release
 86d6c56638539f855cc7703434ff89bce7247d978bb4eab522a9bae9e66bd748 278018 contrib/source/Sources
 354cd4131c98e761a38766488a67bc9acffefe1da352d561d8eb7f9767dedfe6 84833 contrib/source/Sources.gz
 00453a3ffcef0a5dc26a03e61ee9bdfa657444f490c9c7f919939d89b1276d23 148 contrib/source/Release
 2659fe261aac0cfe6ebe251aeb5b74e7c125a63716c086c918be18f9e1e1ec5c 709369 non-free/binary-i386/Packages
 d2e09634a81f7a8f8062344e8e1fcfdd71225e450b068716f48938d2f5c97ed5 175701 non-free/binary-i386/Packages.gz
 702113b6e7b1f61900e5716bf2bf49016c64d6377121d5391b3a6f45a77ae7d7 147 non-free/binary-i386/Release
 291ff577846600849c2acd65082b764033a70a40758cd37b48f1b09e866bc1b4 719937 non-free/binary-amd64/Packages
 6ebc1eaea402a0b02cc963e478ffdc76289e7808a5c87de78f92875b7c0e7ec9 179794 non-free/binary-amd64/Packages.gz
 f2132e0e624aa3889e35c18c41f497f47dfde82d938349cedd0bb489f66c085a 148 non-free/binary-amd64/Release
 b27503ebb08a72fda03a10cc6b44c30a7e0b17a0f41a39b7eb44c367eb6c5490 468876 non-free/source/Sources
 d874ea6ef99713384f9dbb729a7e40d41458b481f52c9c2395835fcfa70b4005 137694 non-free/source/Sources.gz
 2f62fb0934fd5b5581c96e20b38fe796542d766245506f6c917d2aead048247e 149 non-free/source/Release
`
