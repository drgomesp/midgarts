package grf

import (
	"bytes"
	"compress/zlib"
	"io/ioutil"
)

func decompress(data []byte) ([]byte, error) {
	zlibReader, err := zlib.NewReader(bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	defer zlibReader.Close()

	data, err = ioutil.ReadAll(zlibReader)
	if err != nil {
		return nil, err
	}

	return data, nil
}
