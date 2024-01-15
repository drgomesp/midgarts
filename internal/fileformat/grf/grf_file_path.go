package grf

import (
	"path/filepath"
	"strings"

	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/korean"
)

type Path struct {
	bytes       []byte
	windows1252 string
	korean      string
}

func NewFilePath(b []byte) (*Path, error) {
	p := Path{
		bytes: b,
	}

	n, err := charmap.Windows1252.NewDecoder().Bytes(b)
	if err != nil {
		return nil, err
	}
	p.windows1252 = unixifyPaths(string(n))

	k, err := korean.EUCKR.NewDecoder().Bytes(b)
	if err != nil {
		return nil, err
	}

	p.korean = unixifyPaths(string(k))

	return &p, nil
}

// String returns the Windows1252 encoded filepath
func (p Path) String() string {
	return p.windows1252
}

// Korean returns the Korean EUCKR encoded filepath
func (p Path) Korean() string {
	return p.korean
}

// Bytes returns the raw grf extracted byte slice of the filepath
func (p Path) Bytes() []byte {
	return p.bytes
}

// Dir returns the Windows1252 directory of the file path
func (p Path) Dir() string {
	dir, _ := filepath.Split(p.windows1252)
	dir = strings.TrimSuffix(dir, `/`)
	return dir
}

func unixifyPaths(s string) string {
	s = strings.ReplaceAll(s, `\`, `/`)
	s = strings.ToLower(s)
	return s
}
