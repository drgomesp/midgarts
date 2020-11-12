package grf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/pkg/errors"
)

const (
	fileHeaderLength    = 46
	fileHeaderSignature = "Master of Magic"
)

type File struct {
	Header struct {
		Signature       [15]byte
		EncryptionKey   [15]byte
		FileTableOffset uint32
		EntryCount      uint32
		ReservedFiles   uint32
		Version         uint32
	}

	entries map[string]*Entry
	file    *os.File
}

func NewFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("error while opening file: ", err)
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	grfFile := &File{file: f}

	err = grfFile.parseHeader(f, fi)
	if err != nil {
		return nil, errors.Wrap(err, "could not read header")
	}

	err = grfFile.parseEntries(f)
	if err != nil {
		return nil, errors.Wrap(err, "could not read entries")
	}

	return grfFile, nil
}

func (f *File) GetEntries() map[string]*Entry {
	return f.entries
}

func (f *File) GetEntry(name string) (entry *Entry, err error) {
	var exists bool
	if entry, exists = f.entries[name]; !exists {
		return entry, fmt.Errorf("could not find entry '%s'", name)
	}

	_, err = f.file.Seek(int64(entry.Header.Offset)+fileHeaderLength, io.SeekStart)
	if err != nil {
		return nil, err
	}

	data := readNextBytes(f.file, int(entry.Header.CompressedSizeAligned))
	if err = entry.Decode(data); err != nil {
		return nil, err
	}

	return
}

func (f *File) Close() error {
	return f.file.Close()
}

func (f *File) parseHeader(file *os.File, fi os.FileInfo) error {
	err := binary.Read(file, binary.LittleEndian, &f.Header)
	if err != nil {
		return errors.Wrap(err, "could not read file")
	}

	if string(f.Header.Signature[:]) != fileHeaderSignature {
		return errors.New("invalid file signature")
	}

	if f.Header.Version != 0x200 {
		return errors.New("unsupported file version")
	}

	f.Header.FileTableOffset += fileHeaderLength

	if f.Header.FileTableOffset > uint32(fi.Size()) || f.Header.FileTableOffset < 0 {
		return errors.New("invalid file table offset")
	}

	f.Header.EntryCount = f.Header.ReservedFiles - f.Header.EntryCount - 7
	f.entries = make(map[string]*Entry, f.Header.EntryCount)

	return nil
}

func (f *File) parseEntries(file *os.File) error {
	_, _ = file.Seek(int64(f.Header.FileTableOffset), io.SeekStart)

	var compressedSize, uncompressedSize uint32

	_ = binary.Read(file, binary.LittleEndian, &compressedSize)
	_ = binary.Read(file, binary.LittleEndian, &uncompressedSize)

	data, err := decompress(readNextBytes(file, int(compressedSize)))
	if err != nil {
		return err
	}

	for i, offset := 0, 0; i < int(f.Header.EntryCount); i++ {
		var (
			fileName    string
			currentChar byte
			buf         = bytes.NewBufferString(fileName)
		)

		for {
			currentChar = data[offset]
			offset++

			if currentChar == 0 {
				break
			}

			buf.WriteByte(currentChar)
		}

		entry := &Entry{Data: new(bytes.Buffer)}
		fileName = buf.String()

		if err := binary.Read(
			bytes.NewReader(data[offset:offset+entryHeaderLength]),
			binary.LittleEndian, &entry.Header,
		); err != nil {
			return errors.Wrap(err, "could not read file entry header")
		}

		offset += entryHeaderLength

		if entry.Header.Flags&typeFile == 0 {
			continue
		}

		f.entries[fileName] = entry
	}

	return nil
}

func readNextBytes(reader io.Reader, number int) []byte {
	bytesRead := make([]byte, number)

	_, err := reader.Read(bytesRead)
	if err != nil {
		log.Fatal(errors.Wrap(err, "could not read next bytes"))
	}

	return bytesRead
}
