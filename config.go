package dpkg

import (
	"path"
)

var Debug = false
var Strict = true

const ReleaseFileName = "Release"

const DBIndexName = "index.dat"

func DBName(arch Architecture) string { return string(arch) + ".dat" }

func buildDBPath(dataDir string, codeName string, name ...string) string {
	return path.Join(append([]string{dataDir, codeName}, name...)...)
}
