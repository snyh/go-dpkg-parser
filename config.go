package dpkg

import (
	"fmt"
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

type NotFoundError struct {
	resource string
}

func (e NotFoundError) Error() string { return "Not Found resource of " + e.resource }

type FormatError struct {
	t     string
	raw   string
	chain error
}

func (e FormatError) Error() string {
	if e.chain != nil {
		ef, ok := e.chain.(FormatError)
		if ok {
			return fmt.Sprintf("Parsing %q to %q failed at %q", e.raw, e.t+"."+ef.t, ef.raw)
		} else {
			return fmt.Sprintf("Parsing %q to %q failed: %q", e.raw, e.t, e.chain)
		}
	} else {
		return fmt.Sprintf("Parsing %q to %q failed.", e.raw, e.t)
	}
}
