package dpkg

import (
	"fmt"
	"os"
)

// see also deb-src-control(5)

type PackageDatabase struct {
	SourcePackages map[string]ControlFile

	binaryCache map[string]string
	dependCache map[string][]string
}

func NewPackageDatabase(fPath string, keyId string) (*PackageDatabase, error) {
	f, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}

	cfs, err := ParseControlFileGroup(f)
	if err != nil {
		return nil, err
	}
	pd := &PackageDatabase{
		SourcePackages: make(map[string]ControlFile),
		dependCache:    make(map[string][]string),
		binaryCache:    make(map[string]string),
	}
	for _, cf := range cfs {
		key := cf.GetString(keyId)
		pd.SourcePackages[key] = cf

		for _, binary := range cf.GetArrayString("Binary", ",") {
			pd.binaryCache[binary] = key
		}
	}
	fmt.Println("ENEN:", len(pd.SourcePackages))
	return pd, nil
}

func (db *PackageDatabase) Filter(keyId string, value string) []string {
	var ret []string
	for cfName, cf := range db.SourcePackages {
		if cf.GetString(keyId) == value {
			ret = append(ret, cfName)
		}
	}
	return ret
}

func (db *PackageDatabase) FindSourcePackage(binaryName string) (ControlFile, bool) {
	sourcePkg, ok := db.binaryCache[binaryName]
	if !ok {
		fmt.Printf("Warning: Can't find %q's source package\n", binaryName)
	}
	cf, ok := db.SourcePackages[sourcePkg]
	if ok {
		return cf, ok
	}
	return cf, ok
}

func (db *PackageDatabase) QueryBuildDepends(pkgId string) ([]string, []string) {
	cf, ok := db.FindSourcePackage(pkgId)
	if !ok {
		return nil, []string{pkgId}
	}
	cache := make(map[string]bool)

	ret := queryBuildDepends(cf, db, cache)

	var unknown []string
	for id, has := range cache {
		if !has {
			unknown = append(unknown, id)
		}
	}
	return ret, unknown
}

func queryBuildDepends(cf ControlFile, db *PackageDatabase, cache map[string]bool) []string {
	var ret []string
	for _, selfDep := range cf.GetArrayString("Build-Depends", ",") {
		subcf, ok := db.FindSourcePackage(selfDep)
		if !ok {
			cache[selfDep] = false
			continue
		}

		if _, ok := cache[selfDep]; ok {
			continue
		} else {
			cache[selfDep] = true
		}

		for _, subDep := range queryBuildDepends(subcf, db, cache) {
			if _, ok := cache[subDep]; ok {
				continue
			} else {
				cache[subDep] = true
			}
			ret = append(ret, subDep)
		}

		ret = append(ret, selfDep)
	}
	return ret
}
