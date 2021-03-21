package grf

import (
	"bytes"
	"compress/zlib"
	"io"
)

func decompress(data []byte) ([]byte, error) {
	out := new(bytes.Buffer)

	zlibReader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}

	_, err = io.Copy(out, zlibReader)
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
