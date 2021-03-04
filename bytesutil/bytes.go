package bytesutil

import (
	"io"
	"io/ioutil"

	"github.com/pkg/errors"
)

func SkipBytes(buf io.ReadSeeker, n int64) error {
	if _, err := ioutil.ReadAll(io.LimitReader(buf, n)); err != nil {
		return errors.Wrapf(err, "could not skip %d bytes\n", n)
	}

	return nil
}
