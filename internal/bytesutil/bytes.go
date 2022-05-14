package bytesutil

import (
	"encoding/binary"
	"github.com/pkg/errors"
	"io"
	"io/ioutil"
	"strings"
)

func SkipBytes(buf io.ReadSeeker, n int64) error {
	if _, err := ioutil.ReadAll(io.LimitReader(buf, n)); err != nil {
		return errors.Wrapf(err, "could not skip %d bytes\n", n)
	}

	return nil
}

func ReadString(buf io.Reader, length int) (string, error) {
	str := make([]byte, length)
	_ = binary.Read(buf, binary.LittleEndian, &str)

	return strings.Split(string(str), string('\x00'))[0], nil
}
