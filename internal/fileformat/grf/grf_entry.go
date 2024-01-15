package grf

import (
	"github.com/pkg/errors"
	"github.com/project-midgard/midgarts/internal/fileformat/grf/des"
)

type entryFlags byte

const (
	entryHeaderLength = 4 + 4 + 4 + 1 + 4

	entryType              entryFlags = 0x01
	entryTypeEncryptMixed             = 0x02
	entryTypeEncryptHeader            = 0x04
)

// EntryHeader ...
type EntryHeader struct {
	CompressedSize        uint32
	CompressedSizeAligned uint32
	UncompressedSize      uint32
	Flags                 entryFlags
	Offset                uint32
}

// Entry ...
type Entry struct {
	Name *Path

	Header EntryHeader
	Data   []byte
}

// Decode ...
func (e *Entry) Decode(data []byte) error {
	if e.Header.Flags&entryTypeEncryptMixed != 0 {
		des.DecodeFull(data, int(e.Header.CompressedSizeAligned), int(e.Header.CompressedSize))
	} else if e.Header.Flags&entryTypeEncryptHeader != 0 {
		des.DecodeHeader(data)
	}

	if e.Header.CompressedSize == e.Header.UncompressedSize {
		e.Data = data
		return nil
	}

	data, err := decompress(data)
	if err != nil {
		return errors.Wrap(err, "could not decompress entry data")
	}
	e.Data = data

	return nil
}
