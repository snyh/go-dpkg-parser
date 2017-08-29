package dpkg

import (
	"bytes"
	"fmt"
)

type Repository struct {
	dataDir string

	Suites []ReleaseFile

	archives map[string]Archive
	contents map[string]map[string][]string
}

func (r *Repository) AddSuite(host string, codename string, hash string) error {
	rf, err := DownloadReleaseFile(host, codename)
	if err != nil {
		return err
	}
	r.Suites = append(r.Suites, rf)
	return nil
}

func (r *Repository) Archive(arch string) (Archive, error) {
	var err error
	_, ok := r.archives[arch]
	if !ok {
		var indices []IndicesFile
		for _, s := range r.Suites {
			indices = append(indices, s.IndicesFiles(arch, tCONTROLFILES)...)
		}
		r.archives[arch], err = loadOrBuildArchive(r.dataDir, arch, indices)
	}
	return r.archives[arch], err
}

func (r *Repository) Contents(arch string) (map[string][]string, error) {
	var err error
	_, ok := r.contents[arch]
	if !ok {
		var indices []IndicesFile
		for _, s := range r.Suites {
			indices = append(indices, s.IndicesFiles(arch, tCONTENTS)...)
		}
		r.contents[arch], err = loadOrBuildContent(r.dataDir, indices)
	}
	return r.contents[arch], err
}

func NewRepository(dataDir string) *Repository {
	r := &Repository{
		dataDir:  dataDir,
		archives: make(map[string]Archive),
		contents: make(map[string]map[string][]string),
	}
	return r
}

func (r *Repository) Timestamp() string {
	if len(r.Suites) == 0 {
		return "Unknown"
	}
	s := r.Suites[0]
	return s.Date.UTC().String()
}

func (r *Repository) Hash() string {
	var hs []string
	for _, s := range r.Suites {
		hs = append(hs, s.Hash)
	}
	return HashArrayString(hs)
}

func DownloadReleaseFile(repoURL string, suiteName string) (ReleaseFile, error) {
	url := fmt.Sprintf("%s/dists/%s/%s", repoURL, suiteName, ReleaseFileName)
	buf := bytes.NewBuffer(nil)
	err := DownloadTo(url, buf)
	if err != nil {
		return ReleaseFile{}, err
	}
	cf, err := NewControlFile(string(buf.Bytes()))
	if err != nil {
		return ReleaseFile{}, err
	}
	return cf.ToReleaseFile(repoURL)
}
