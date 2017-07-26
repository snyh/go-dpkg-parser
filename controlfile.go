package dpkg

import (
	"bufio"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"strings"
)

type ControlFile map[string]string

func splitControlFileLine(data []byte, noMoreData bool) (int, []byte, error) {
	l := len(data)
	maybeWrap := 0

	EndOfChunk := func(i int) bool { return i+1 == l }
	NextIsNotSpace := func(i int) bool { c := data[i+1]; return c != ' ' && c != '\t' }

	for i, c := range data {
		if EndOfChunk(i) {
			continue
		}
		switch c {
		case '\n':
			if maybeWrap != 0 {
				continue
			} else if NextIsNotSpace(i) {
				return i + 1, data[:i], nil
			}
		case '\\':
			//TODO: strip the wrap characters before returning
			maybeWrap = i
		default:
			if NextIsNotSpace(i) {
				maybeWrap = 0
			}
		}
	}

	if !noMoreData {
		return 0, nil, nil
	} else if l != 0 {
		return l, data, nil
	}

	return 0, nil, fmt.Errorf("End of file")

}

func NewControlFile(r io.Reader) (ControlFile, error) {
	s := bufio.NewScanner(r)
	s.Split(splitControlFileLine)

	f := make(ControlFile)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		d := strings.SplitN(line, ":", 2)
		if len(d) != 2 {
			if Strict {
				return nil, fmt.Errorf("NewControlFile there has %d separators at:%q", len(d), line)
			}
			continue
		}
		f[strings.ToLower(d[0])] = strings.Trim(d[1], " \n")
	}
	return f, nil
}

func (d ControlFile) GetString(key string) string {
	return strings.TrimSpace(d[strings.ToLower(key)])
}

func (d ControlFile) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	for k, v := range d {
		buf.WriteString(fmt.Sprintf("%s : %s\n", k, v))
	}
	return buf.Bytes()
}

func (d ControlFile) GetArrayString(key string, sep string) []string {
	s := d.GetString(key)
	return getArrayString(s, sep)
}

func (d ControlFile) GetMultiline(key string) []string {
	return getMultiline(d.GetString(key))
}

func LoadControlFileGroup(fPath string) ([]ControlFile, error) {
	f, err := os.Open(fPath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	if strings.HasSuffix(strings.ToLower(fPath), ".gz") {
		gr, err := gzip.NewReader(f)
		if err != nil {
			return nil, fmt.Errorf("parsePackageDBComponent handle zip file(%v) error:%v", fPath, err)
		}
		defer gr.Close()
		return ParseControlFileGroup(gr)
	}

	return ParseControlFileGroup(f)
}

func ParseControlFileGroup(r io.Reader) ([]ControlFile, error) {
	s := bufio.NewScanner(r)
	s.Buffer(nil, 512*1024)
	splitFn := func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		l := len(data)
		for i, c := range data {
			if c == '\n' {
				if i+1 < l && data[i+1] == '\n' {
					return i + 2, data[:i], nil
				}
				if i+1 == l && atEOF {
					return i + 1, data[:i], nil
				}
			}
		}
		if !atEOF {
			return 0, nil, nil
		}

		if atEOF && l != 0 {
			return l, data, nil
		}

		return l, data, fmt.Errorf("end of file")
	}

	s.Split(splitFn)

	var ts []ControlFile
	for s.Scan() {
		cf, err := NewControlFile(bytes.NewBuffer(s.Bytes()))
		if err != nil {
			return nil, err
		}
		ts = append(ts, cf)
	}
	return ts, nil
}

func getArrayString(s string, sep string) []string {
	var r []string
	for _, c := range strings.Split(s, sep) {
		r = append(r, strings.TrimSpace(c))
	}
	return r
}

func getMultiline(s string) []string {
	if s == "" {
		return nil
	}
	var ret []string

	for _, f := range strings.Split(s, "\n") {
		ret = append(ret, strings.TrimSpace(f))
	}
	return ret

}
