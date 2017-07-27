package dpkg

import (
	"fmt"
	"os"
)

var Debug = false

func DebugPrintf(fmtStr string, args ...interface{}) {
	if Debug {
		fmt.Fprintf(os.Stderr, fmtStr, args...)
	}
}
func DebugPrintln(args ...interface{}) {
	if Debug {
		fmt.Fprintln(os.Stderr, args...)
	}
}

var Strict = true

const ScanBufferSize = 512 * 1024

const ReleaseFileName = "Release"

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
