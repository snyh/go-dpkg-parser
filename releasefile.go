package dpkg

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"
)

type PackagesFileInfo struct {
	Size         uint64
	Path         string
	Gzip         bool
	MD5          string
	Architecture Architecture
}

type ReleaseFile struct {
	Date          string
	CodeName      string
	Description   string
	Components    []string
	Architectures Architectures
	fileInfos     []PackagesFileInfo
}

func NewReleaseFile(r io.Reader) (ReleaseFile, error) {
	cf, err := NewControlFile(r)
	if err != nil {
		return ReleaseFile{}, err
	}
	return cf.ToReleaseFile()
}

// GetReleaseFile load ReleaseFile from dataDir with codeName
func GetReleaseFile(dataDir string, codeName string) (ReleaseFile, error) {
	f, err := os.Open(buildDBPath(dataDir, codeName, ReleaseFileName))
	if err != nil {
		return ReleaseFile{}, fmt.Errorf("GetReleaseFile open file error: %v", err)
	}
	cf, err := NewControlFile(f)
	if err != nil {
		return ReleaseFile{}, err
	}
	return cf.ToReleaseFile()
}

// ToReleaseFile build a new ReleaseFile by reading contents from r
func (cf ControlFile) ToReleaseFile() (ReleaseFile, error) {
	rf := ReleaseFile{}

	for _, arch := range cf.GetArrayString("architectures", " ") {
		rf.Architectures = append(rf.Architectures, Architecture(arch))
	}
	rf.Date = cf.GetString("date")

	rf.Date = cf.GetString("date")
	rf.CodeName = cf.GetString("codename")
	rf.Description = cf.GetString("description")
	rf.Date = cf.GetString("date")
	rf.Components = cf.GetArrayString("components", " ")

	var ps []PackagesFileInfo
	for _, v := range cf.GetMultiline("md5sum") {
		fs := strings.Split(strings.TrimSpace(v), " ")
		if len(fs) != 3 {
			continue
		}
		size, err := strconv.Atoi(fs[1])
		if err != nil {
			continue
		}

		ps = append(ps, PackagesFileInfo{
			Size: uint64(size),
			Path: fs[2],
			Gzip: strings.HasSuffix(fs[2], ".gz"),
			MD5:  fs[0],
		})
	}
	rf.fileInfos = ps
	return rf, rf.valid()
}
func (rf ReleaseFile) valid() error {
	if rf.CodeName == "" {
		return fmt.Errorf("NewReleaseFile input data is invalid. Without codename")
	}
	if len(rf.Components) == 0 {
		return fmt.Errorf("NewReleaseFile input data is invalid. Without any components")
	}

	if len(rf.FileInfos()) == 0 {
		return fmt.Errorf("NewReleaseFile input data is invalid. Without any valid file")
	}
	return nil
}

func (rf ReleaseFile) Hash() string {
	var data []byte
	for _, finfo := range rf.FileInfos() {
		data = append(data, ([]byte)(finfo.MD5)...)
	}
	return fmt.Sprintf("%x", md5.Sum(data))
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
	zip := raw + ".gz"
	found := false
	var b PackagesFileInfo
	for _, f := range rf.fileInfos {
		if f.Path != raw && f.Path != zip {
			continue
		}
		found = true
		b = f
		if f.Gzip {
			return f, found
		}
	}
	return b, found
}
func (rf ReleaseFile) FileInfos() []PackagesFileInfo {
	var set = make(map[string]PackagesFileInfo)
	for _, component := range rf.Components {
		for _, arch := range rf.Architectures {
			raw := component + "/binary-" + string(arch) + "/Packages"
			p, ok := rf.findComponent(raw)
			if ok {
				p.Architecture = arch
				set[raw] = p
			}
		}

		raw := component + "/source/Sources"
		if p, ok := rf.findComponent(raw); ok {
			p.Architecture = "source"
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
