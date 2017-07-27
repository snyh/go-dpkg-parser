package dpkg

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type ControlFile map[string]string

func NewControlFiles(r io.Reader) ([]ControlFile, error) {
	s := bufio.NewScanner(r)
	s.Buffer(nil, ScanBufferSize)

	s.Split(splitControlFile)

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

func NewControlFile(r io.Reader) (ControlFile, error) {
	s := bufio.NewScanner(r)
	s.Split(splitControlFileLine)

	f := make(ControlFile)
	for s.Scan() {
		line := s.Text()
		if line == "" {
			continue
		}
		kv := strings.SplitN(line, ":", 2)
		if len(kv) != 2 {
			if Strict {
				return nil, FormatError{"ConrolfileField", line, nil}
			}
			continue
		}
		f[strings.ToLower(kv[0])] = strings.TrimRight(kv[1], " \n")
	}
	if len(f) == 0 {
		return nil, NotFoundError{"Empty result"}
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

func splitControlFileLine(chunk []byte, noMoreChunk bool) (int, []byte, error) {
	chunkSize := len(chunk)
	maybeWrap := 0

	EndOfChunk := func(i int) bool { return i+1 == chunkSize }
	NextIsNotSpace := func(i int) bool { c := chunk[i+1]; return c != ' ' && c != '\t' }

	for i, c := range chunk {
		if EndOfChunk(i) {
			continue
		}
		switch c {
		case '\n':
			if maybeWrap != 0 {
				continue
			} else if NextIsNotSpace(i) {
				return i + 1, chunk[:i], nil
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

	if !noMoreChunk {
		return 0, nil, nil
	} else if chunkSize > 0 {
		return chunkSize, chunk, nil
	}

	return 0, nil, fmt.Errorf("End of file")

}
func splitControlFile(chunk []byte, noMoreChunk bool) (int, []byte, error) {
	chunkSize := len(chunk)

	EndOfChunk := func(i int) bool { return i+1 == chunkSize }

	for i, c := range chunk {
		if c != '\n' {
			continue
		}

		if !EndOfChunk(i) && chunk[i+1] == '\n' {
			return i + 2, chunk[:i], nil
		}
		if EndOfChunk(i) && noMoreChunk {
			return i + 1, chunk[:i], nil
		}
	}

	if !noMoreChunk {
		return 0, nil, nil
	} else if chunkSize > 0 {
		return chunkSize, chunk, nil
	}
	return 0, nil, fmt.Errorf("end of file")
}
