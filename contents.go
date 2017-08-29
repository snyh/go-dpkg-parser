package dpkg

import (
	"bufio"
	"strings"
)

func parsePkgNameInContent(str string) string {
	return last(strings.Split(str, "/"))
}
func parseContentIndices(f string) map[string][]string {
	r, err := ReadFile(f)
	if err != nil {
		return nil
	}
	defer r.Close()

	ret := make(map[string][]string)
	buf := bufio.NewReader(r)
	for {
		line, err := buf.ReadString('\n')
		if err != nil {
			break
		}
		if line == "" {
			continue
		}

		fs := getArrayString(line, " ")
		var fpath, pkgname string
		switch n := len(fs); n {
		case 1, 0:
			DebugPrintf("Invalid content indices line: %q\n", line)
			continue
		case 2:
			fpath = fs[0]
			pkgname = parsePkgNameInContent(fs[1])
		default:
			fpath = strings.Join(fs[:n-1], " ")
			pkgname = parsePkgNameInContent(fs[n-1])
		}
		if pkgname == "" || fpath == "" {
			panic("Internal error" + line)
		}
		ret[pkgname] = append(ret[pkgname], fpath)
	}
	return ret
}
