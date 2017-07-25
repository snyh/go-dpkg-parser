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

func NewControlFile(r io.Reader) (ControlFile, error) {
	splitFn := func(data []byte, atEOF bool) (advance int, toke []byte, err error) {
		l := len(data)
		for i, c := range data {
			if c == '\n' {
				if i+1 < l && (data[i+1] != ' ' && data[i+1] != '\t') {
					return i + 1, data[:i], nil
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

		return l, data, fmt.Errorf("End of file")
	}

	s := bufio.NewScanner(r)
	s.Split(splitFn)

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
	return d[strings.ToLower(key)]
}

func (d ControlFile) Bytes() []byte {
	buf := bytes.NewBuffer(nil)
	for k, v := range d {
		buf.WriteString(fmt.Sprintf("%s : %s\n", k, v))
	}
	return buf.Bytes()
}

func (d ControlFile) GetArrayString(key string, sep string) []string {
	var r []string
	for _, c := range strings.Split(d.GetString(key), sep) {
		r = append(r, strings.TrimSpace(c))
	}
	return r
}

func (d ControlFile) GetMultiline(key string) []string {
	return strings.Split(d[strings.ToLower(key)], "\n")
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
