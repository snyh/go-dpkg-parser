package dpkg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

var Strict = true

type ControlFile map[string]string

func NewControlFile(data []byte) (ControlFile, error) {
	r := bytes.NewBuffer(data)
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

	f := make(ControlFile)
	s := bufio.NewScanner(r)
	s.Split(splitFn)
	for s.Scan() {
		line := s.Text()
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

func (d ControlFile) GetArrayString(key string) []string {
	var r []string
	for _, c := range strings.Split(d.GetString(key), " ") {
		r = append(r, c)
	}
	return r
}

func (d ControlFile) GetMultiline(key string) []string {
	return strings.Split(d[strings.ToLower(key)], "\n")
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
		cf, err := NewControlFile(s.Bytes())
		if err != nil {
			return nil, err
		}
		ts = append(ts, cf)
	}
	return ts, nil
}
