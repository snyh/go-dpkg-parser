package dpkg

import (
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
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

const DBIndexName = "index.dat"
const ReleaseFileName = "Release"

func DBName(arch Architecture) string { return string(arch) + ".dat" }

type ReleaseFile struct {
	Date          string
	CodeName      string
	Description   string
	Components    []string
	Architectures Architectures
	fileInfos     []PackagesFileInfo
}

// GetReleaseFile load ReleaseFile from dataDir with codeName
func GetReleaseFile(dataDir string, codeName string) (ReleaseFile, error) {
	bs, err := ioutil.ReadFile(buildDBPath(dataDir, codeName, ReleaseFileName))
	if err != nil {
		return ReleaseFile{}, fmt.Errorf("GetReleaseFile open file error: %v", err)
	}
	cf, err := NewControlFile(bs)
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
	if rf.CodeName == "" || len(rf.FileInfos()) == 0 || len(rf.Components) == 0 {
		return fmt.Errorf("NewReleaseFile input data is invalid. %v", rf)
	}
	return nil
}

func buildDBPath(dataDir string, codeName string, name ...string) string {
	return path.Join(append([]string{dataDir, codeName}, name...)...)
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

func (rf ReleaseFile) FileInfos() []PackagesFileInfo {
	var set = make(map[string]PackagesFileInfo)
	for _, arch := range rf.Architectures {
		for _, component := range rf.Components {
			raw := component + "/binary-" + string(arch) + "/Packages"
			zip := raw + ".gz"
			for _, f := range rf.fileInfos {
				if f.Path != raw && f.Path != zip {
					continue
				}
				_, ok := set[raw]
				if !ok {
					//store it if there hasn't content
					f.Architecture = arch
					set[raw] = f
				}
				if f.Gzip {
					//overwrite if it support gzip
					f.Architecture = arch
					set[raw] = f
				}
			}
		}
	}

	var r = make(PackagesFileInfos, 0)
	for _, f := range set {
		r = append(r, f)
	}
	sort.Sort(r)
	return r
}

func HashFile(fpath string) string {
	f, err := os.Open(fpath)
	if err != nil {
		return ""
	}
	defer f.Close()

	hash := md5.New()
	_, err = io.Copy(hash, f)
	if err != nil {
		return ""
	}
	var r [16]byte
	copy(r[:], hash.Sum(nil))
	return fmt.Sprintf("%x", r)
}
