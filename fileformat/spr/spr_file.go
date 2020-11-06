package spr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strconv"

	"github.com/pkg/errors"
)

const (
	HeaderSignature = "SP"
	PaletteSize     = 256 * 4
)

type SpriteFile struct {
	Header struct {
		Signature string
		Version   float32

		IndexedFrameCount uint16
		RGBAFrameCount    uint16
		RGBAIndex         uint16
	}

	Palette *bytes.Buffer
}

func Load(buf io.Reader) (file *SpriteFile, err error) {
	file = new(SpriteFile)

	if err := file.parseHeader(buf); err != nil {
		return nil, err
	}

	return file, nil
}

func (f *SpriteFile) parseHeader(buf io.Reader) error {
	var signature [2]byte
	_ = binary.Read(buf, binary.LittleEndian, &signature)

	signatureStr := string(signature[:])
	if signatureStr != HeaderSignature {
		return fmt.Errorf("invalid signature: %s\n", signature)
	}

	var major, minor byte
	_ = binary.Read(buf, binary.LittleEndian, &major)
	_ = binary.Read(buf, binary.LittleEndian, &minor)

	version, err := strconv.ParseFloat(fmt.Sprintf("%d.%d", major, minor), 32)
	if err != nil {
		return errors.Wrapf(err, "invalid version: %s\n", strconv.FormatFloat(version, 'E', -1, 32))
	}

	var indexedFrameCount, rgbaFrameCount uint16
	if err = binary.Read(buf, binary.LittleEndian, &indexedFrameCount); err != nil {
		return errors.Wrap(err, "could not read indexed frame count")
	}

	if version > 1.1 {
		if err = binary.Read(buf, binary.LittleEndian, &rgbaFrameCount); err != nil {
			return errors.Wrap(err, "could not read rgba frame count")
		}
	}

	f.Header.Signature = signatureStr
	f.Header.Version = float32(version)
	f.Header.IndexedFrameCount = indexedFrameCount
	f.Header.RGBAFrameCount = rgbaFrameCount
	f.Header.RGBAIndex = indexedFrameCount
	f.Palette = bytes.NewBuffer(make([]byte, PaletteSize))

	return nil
}
