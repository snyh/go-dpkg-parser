package dpkg

import (
	"os"
	"path"
)

func BuildIndices(rootDir string, buildFn BuildFunc, indices []IndicesFile) (string, error) {
	err := downloadIndicesContent(rootDir, indices)
	if err != nil {
		return "", err
	}
	var fpaths []string
	var hashs []string
	for _, index := range indices {
		fpaths = append(fpaths, path.Join(rootDir, index.CachePath))
		hashs = append(hashs, index.Hash)
	}
	cacheFile := path.Join(rootDir, "db", HashArrayString(hashs))
	err = EnsureDirectory(path.Dir(cacheFile))
	if err != nil {
		return "", err
	}

	_, err = os.Stat(cacheFile)
	if err == nil {
		return cacheFile, nil
	}
	cache, err := buildFn(fpaths...)
	if err != nil {
		return "", err
	}
	return cacheFile, storeGOB(cacheFile, cache)
}

func downloadIndicesContent(rootDir string, indices []IndicesFile) error {
	for _, index := range indices {
		local := path.Join(rootDir, index.CachePath)
		hash, _ := HashFile(local)
		if hash == index.Hash {
			DebugPrintf("%q is cached\n", local)
			continue
		}
		err := DownloadToFile(
			index.Url,
			path.Join(rootDir, index.CachePath),
		)
		if err != nil {
			return err
		}
	}
	return nil
}

type BuildFunc func(files ...string) (interface{}, error)

func loadOrBuildContent(baseDir string, files []IndicesFile) (map[string][]string, error) {
	cache := make(map[string][]string)
	cacheFile, err := BuildIndices(baseDir, BuildContent, files)
	if err != nil {
		return nil, err
	}
	err = loadGOB(cacheFile, &cache)
	return cache, err
}

func loadOrBuildArchive(baseDir string, arch string, files []IndicesFile) (Archive, error) {
	cache := NewArchive(arch)
	cacheFile, err := BuildIndices(baseDir, BuildArchive, files)
	if err != nil {
		return cache, err
	}
	err = loadGOB(cacheFile, &cache)
	return cache, err
}

func BuildArchive(files ...string) (interface{}, error) {
	cache := NewArchive("")
	for _, fpath := range files {
		cfs, err := LoadPackages(fpath)
		if err != nil {
			return cache, err
		}
		for _, cf := range cfs {
			name := cf.Get("Package")
			cache.Packages[name] = cf
			for _, p := range parseProvides(cf.Get("provides")) {
				cache.Virtuals[p] = append(cache.Virtuals[p], name)
			}
		}
	}
	return cache, nil
}

func BuildContent(files ...string) (interface{}, error) {
	ret := make(map[string][]string)
	for _, f := range files {
		for name, ps := range parseContentIndices(f) {
			ret[name] = ps
		}
	}
	return ret, nil
}
