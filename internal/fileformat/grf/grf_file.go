package grf

import (
	"bufio"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/project-midgard/midgarts/internal/fileformat/act"
	"github.com/project-midgard/midgarts/internal/fileformat/spr"

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
		return nil, err
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
		if e.Name.String() == name {
			entry = e
			break
		}
	}

	if entry == nil {
		return nil, fmt.Errorf("could not find entry '%s'", name)
	}

	if len(entry.Data) != 0 {
		return entry, nil
	}

	_, err = f.file.Seek(int64(entry.Header.Offset)+fileHeaderLength, io.SeekStart)
	if err != nil {
		return entry, err
	}

	data, err := readNextBytes(f.file, int(entry.Header.CompressedSizeAligned))
	if err != nil {
		return entry, err
	}

	if err = entry.Decode(data); err != nil {
		return entry, err
	}

	return
}

type ActionSpriteFilePair struct {
	ACT *act.ActionFile
	SPR *spr.SpriteFile
}

func (f *File) GetSpriteFiles(name string) (ActionSpriteFilePair, error) {
	e, err := f.GetEntry(fmt.Sprintf("%s.act", name))
	if err != nil {
		return ActionSpriteFilePair{}, err
	}

	actFile, err := act.Load(e.Data)
	if err != nil {
		return ActionSpriteFilePair{}, err
	}

	e, err = f.GetEntry(fmt.Sprintf("%s.spr", name))
	if err != nil {
		return ActionSpriteFilePair{}, err
	}

	sprFile, err := spr.Load(e.Data)
	if err != nil {
		return ActionSpriteFilePair{}, err
	}

	return ActionSpriteFilePair{
		ACT: actFile,
		SPR: sprFile,
	}, nil
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

	zlibReader, err := zlib.NewReader(file)
	if err != nil {
		return errors.Wrap(err, "could instantiate zlib reader")
	}
	var (
		reader     = bufio.NewReader(zlibReader)
		uniqueDirs = make(map[string]bool)
	)

	for i := 0; i < int(f.Header.EntryCount); i++ {
		fileNameBytes, err := reader.ReadBytes(0)
		if err != nil {
			return errors.Wrap(err, "could not parse entry file name")
		}

		entryPath, err := NewFilePath(fileNameBytes[0 : len(fileNameBytes)-1])
		if err != nil {
			return errors.Wrap(err, "decoding entry path")
		}

		entry := &Entry{Data: []byte{}}

		if err = binary.Read(reader, binary.LittleEndian, &entry.Header); err != nil {
			return errors.Wrap(err, "could not read file entry header")
		}

		if entry.Header.Flags&entryType == 0 {
			continue
		}

		entry.Name = entryPath
		dir := entryPath.Dir()
		uniqueDirs[dir] = true

		f.entries[dir] = append(f.entries[dir], entry)
	}

	var dirs []string
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
					return errors.Wrap(err, "could not insert tree nodes")
				}
			}
		} else {
			if err = f.entriesTree.Insert(dir, f.entries[dir]); err != nil {
				return errors.Wrap(err, "could not insert tree nodes")
			}
		}
	}

	return nil
}

func readNextBytes(reader io.Reader, number int) ([]byte, error) {
	bytesRead := make([]byte, number)

	n, err := reader.Read(bytesRead)
	if err != nil {
		return nil, errors.Wrap(err, "could not read next bytes")
	}
	if n != number {
		return nil, errors.Wrapf(err, "could not read next bytes: want %d, got %d", number, n)
	}
	return bytesRead, nil
}
