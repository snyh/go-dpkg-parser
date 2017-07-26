package dpkg

import (
	"sort"
)

type PackageDB map[string]BinaryPackage

type PackageDBIndex struct {
	PackagePaths map[Architecture]map[string]string
	DBPaths      map[Architecture]string
}

func NewPackagesaDBIndex(indexFile string) (*PackageDBIndex, error) {
	return loadPackagesaDBIndex(indexFile)
}

func (dbi PackageDBIndex) DBPath(arch Architecture) (string, bool) {
	p, ok := dbi.DBPaths[arch]
	return p, ok
}
func (dbi PackageDBIndex) Architectures() Architectures {
	var archs Architectures
	for arch := range dbi.DBPaths {
		archs = append(archs, arch)
	}
	sort.Sort(archs)
	return archs
}
func (dbi PackageDBIndex) PackageArchitectures(pid string) Architectures {
	var r Architectures
	for arch, paths := range dbi.PackagePaths {
		_, ok := paths[pid]
		if !ok {
			continue
		}
		r = append(r, arch)
	}
	sort.Sort(r)
	return r
}
func (dbi PackageDBIndex) PackagePath(pid string, arch Architecture) (string, bool) {
	ps, ok := dbi.PackagePaths[arch]
	if !ok {
		return "", false
	}
	p, ok := ps[pid]
	return p, ok
}

func createPackageIndex(dbsPath map[Architecture]string, dbs map[Architecture]PackageDB) PackageDBIndex {
	index := PackageDBIndex{
		DBPaths:      dbsPath,
		PackagePaths: make(map[Architecture]map[string]string),
	}

	for arch, db := range dbs {
		index.PackagePaths[arch] = make(map[string]string)
		for _, t := range db {
			index.PackagePaths[arch][t.Package] = t.Filename
		}
	}
	return index
}

func createPackageDB(sourcePaths []string) (PackageDB, error) {
	r := make(map[string]BinaryPackage)
	for _, source := range sourcePaths {
		cs, err := LoadBinaryPackages(source)
		if err != nil {
			return nil, err
		}
		for _, t := range cs {
			r[t.Package] = t
		}
	}
	return r, nil
}
