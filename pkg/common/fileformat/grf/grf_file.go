package grf

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"golang.org/x/text/encoding/charmap"

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

	// Maps to directory -> array of entries
	entries     map[string][]*Entry
	entriesTree *EntryTree

	file *os.File
}

func Load(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		log.Fatal("error while opening file: ", err)
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}

	grfFile := &File{file: f, entriesTree: &EntryTree{}}

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

func (f *File) GetEntryDirectories() map[string][]*Entry {
	return f.entries
}

func (f *File) GetEntries(dir string) []*Entry {
	return f.entries[dir]
}

func (f *File) GetEntryTree() *EntryTree {
	return f.entriesTree
}

func (f *File) GetEntry(name string) (entry *Entry, err error) {
	name = strings.ToLower(name)
	var entries []*Entry
	var exists bool
	dir, _ := filepath.Split(name)
	dir = strings.TrimSuffix(dir, `/`)

	if entries, exists = f.entriesTree.Find(dir); !exists {
		return entry, fmt.Errorf("could not find directory '%s'", dir)
	}

	for _, e := range entries {
		if e.Name == name {
			entry = e
			break
		}
	}

	if entry == nil {
		return nil, fmt.Errorf("could not find entry '%s'", name)
	}

	_, err = f.file.Seek(int64(entry.Header.Offset)+fileHeaderLength, io.SeekStart)
	if err != nil {
		return entry, nil
	}

	data := readNextBytes(f.file, int(entry.Header.CompressedSizeAligned))
	if err = entry.Decode(data); err != nil {
		return entry, err
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

	sig := string(f.Header.Signature[:])
	if sig != fileHeaderSignature {
		return fmt.Errorf("invalid file signature '%s'", sig)
	}

	if f.Header.Version != 0x200 {
		return errors.New("unsupported file version")
	}

	f.Header.FileTableOffset += fileHeaderLength

	if f.Header.FileTableOffset > uint32(fi.Size()) || f.Header.FileTableOffset < 0 {
		return errors.New("invalid file table offset")
	}

	f.Header.EntryCount = f.Header.ReservedFiles - f.Header.EntryCount - 7
	f.entries = make(map[string][]*Entry, f.Header.EntryCount)

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

	var dirs []string
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

		fileName = buf.String()
		var d []byte
		if d, err = charmap.Windows1252.NewDecoder().Bytes([]byte(fileName)); err != nil {
			panic(err)
		}

		entry := &Entry{Data: new(bytes.Buffer)}

		if err = binary.Read(
			bytes.NewReader(data[offset:offset+entryHeaderLength]),
			binary.LittleEndian, &entry.Header,
		); err != nil {
			return errors.Wrap(err, "could not read file entry header")
		}

		offset += entryHeaderLength

		if entry.Header.Flags&entryType == 0 {
			continue
		}

		properFileName := strings.ToLower(strings.ReplaceAll(string(d), `\`, `/`))
		entry.Name = properFileName
		dir, file := filepath.Split(properFileName)
		dir = strings.TrimSuffix(dir, `/`)
		dirs = append(dirs, dir)

		if strings.HasSuffix(file, ".spr") || strings.HasSuffix(file, ".act") {
			f.entries[dir] = append(f.entries[dir], entry)
		}
	}

	uniqueDirs := map[string]byte{}

	for _, dir := range dirs {
		if _, exists := uniqueDirs[dir]; !exists {
			uniqueDirs[dir] = 0
		}
	}

	dirs = []string{}
	for dir := range uniqueDirs {
		dirs = append(dirs, dir)
	}

	sort.Slice(dirs, func(i, j int) bool {
		return strings.ToLower(dirs[i]) < strings.ToLower(dirs[j])
	})

	for _, dir := range dirs {
		if _, exists := f.entriesTree.Find(dir); exists {
			var toInsert []*Entry

			for _, e := range f.entries[dir] {
				toInsert = append(toInsert, e)
			}

			if len(toInsert) > 0 {
				if err = f.entriesTree.Insert(dir, toInsert); err != nil {
					log.Fatalf("could not insert tree nodes: %v", err)
				}
			}
		} else {
			if err = f.entriesTree.Insert(dir, f.entries[dir]); err != nil {
				log.Fatalf("could not insert tree nodes: %v", err)
			}
		}
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
