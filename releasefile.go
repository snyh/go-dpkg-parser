package dpkg

import (
	"fmt"
	"io"
	"io/ioutil"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	tCONTENTS     = "contents"
	tCONTROLFILES = "controlfiles"
)

type PackagesFileInfo struct {
	Size         uint64
	Path         string
	MD5          string
	Architecture string

	Type string
}

type ReleaseFile struct {
	Date          time.Time
	Suite         string
	Description   string
	Components    []string
	Architectures Architectures
	fileInfos     []PackagesFileInfo
	Hash          string
}

func NewReleaseFile(r io.Reader) (ReleaseFile, error) {
	bs, err := ioutil.ReadAll(r)
	if err != nil {
		return ReleaseFile{}, err
	}
	cf, err := NewControlFile(string(bs))
	if err != nil {
		return ReleaseFile{}, err
	}
	return cf.ToReleaseFile()
}

// GetReleaseFile load ReleaseFile from dataDir with suite
func LoadReleaseFile(path string) (ReleaseFile, error) {
	bs, err := ioutil.ReadFile(path)
	if err != nil {
		return ReleaseFile{}, fmt.Errorf("GetReleaseFile open file error: %v", err)
	}

	cf, err := NewControlFile(string(bs))
	if err != nil {
		return ReleaseFile{}, err
	}
	return cf.ToReleaseFile()
}

func parseSuiteDate(str string) (time.Time, error) {
	for _, layout := range []string{time.RFC1123, time.RFC3339} {
		v, err := time.Parse(layout, str)
		if err == nil {
			return v, nil
		}
	}
	return time.Time{}, fmt.Errorf("unknown the layout of date string: %q", str)
}

// ToReleaseFile build a new ReleaseFile by reading contents from r
func (cf ControlFile) ToReleaseFile() (ReleaseFile, error) {
	rf := ReleaseFile{}

	for _, arch := range cf.GetArray("architectures", " ") {
		rf.Architectures = append(rf.Architectures, arch)
	}
	rf.Suite = cf.Get("suite")
	if rf.Suite == "" {
		rf.Suite = cf.Get("codename")
	}
	rf.Description = cf.Get("description")
	var err error
	rf.Date, err = parseSuiteDate(cf.Get("date"))
	if err != nil {
		DebugPrintf("Unknown release date: %s\n", err)
	}

	rf.Components = cf.GetArray("components", " ")
	rf.Hash = HashBytes([]byte(cf.Raw))

	var ps []PackagesFileInfo
	for _, v := range cf.GetMultiline("md5sum") {
		fs := getArrayString(v, " ")
		if len(fs) != 3 {
			DebugPrintf("Ignore:%q %+v (%d)\n", v, fs, len(fs))
			continue
		}
		size, err := strconv.Atoi(fs[1])
		if err != nil {
			DebugPrintf("Components size field invalid %q\n", v)
			continue
		}

		ps = append(ps, PackagesFileInfo{
			Size: uint64(size),
			Path: fs[2],
			MD5:  fs[0],
		})
	}
	rf.fileInfos = ps
	return rf, rf.valid()
}
func (rf ReleaseFile) valid() error {
	if rf.Suite == "" {

		return fmt.Errorf("NewReleaseFile input data is invalid. Without suite name")
	}
	if len(rf.Components) == 0 {
		return fmt.Errorf("NewReleaseFile input data is invalid. Without any components")
	}

	if len(rf.FileInfos()) == 0 {
		return fmt.Errorf("NewReleaseFile input data is invalid. Without any valid fileinfos")
	}
	return nil
}

type PackagesFileInfos []PackagesFileInfo

func (infos PackagesFileInfos) Len() int {
	return len(infos)
}
func (infos PackagesFileInfos) Less(i, j int) bool {
	return infos[i].Path < infos[j].Path
}
func (infos PackagesFileInfos) Swap(i, j int) {
	infos[i], infos[j] = infos[j], infos[i]
}

func (rf ReleaseFile) findComponent(raw string) (PackagesFileInfo, bool) {
	found := false
	var fallback PackagesFileInfo
	for _, f := range rf.fileInfos {
		if f.Path != raw && f.Path != raw+".gz" && f.Path != raw+".bz2" {
			continue
		}
		found, fallback = true, f

		switch strings.ToLower(path.Ext(f.Path)) {
		case ".gz", ".bz2":
			return f, found
		}
	}
	return fallback, found
}
func (rf ReleaseFile) FileInfos() []PackagesFileInfo {
	var set = make(map[string]PackagesFileInfo)
	for _, component := range rf.Components {
		for _, arch := range rf.Architectures {
			raw := component + "/binary-" + string(arch) + "/Packages"
			if p, ok := rf.findComponent(raw); ok {
				p.Architecture = arch
				p.Type = tCONTROLFILES
				set[raw] = p

			}
			raw = component + "/Contents-" + string(arch)
			if p, ok := rf.findComponent(raw); ok {
				p.Architecture = arch
				p.Type = tCONTENTS
				set[raw] = p
			}
		}

		raw := component + "/source/Sources"
		if p, ok := rf.findComponent(raw); ok {
			p.Architecture = "source"
			p.Type = tCONTROLFILES
			set[raw] = p
		}
	}

	var r = make(PackagesFileInfos, 0)
	for _, f := range set {
		r = append(r, f)
	}
	sort.Sort(r)
	return r
}
