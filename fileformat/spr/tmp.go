package spr

//
//import (
//	"bufio"
//	"bytes"
//	"encoding/binary"
//	"fmt"
//	"io"
//	"io/ioutil"
//	"log"
//	"strconv"
//
//	"github.com/pkg/errors"
//
//	bytesUtil "github.com/project-midgard/midgarts/bytes"
//)
//
//type FileType int
//
//const (
//	PaletteSize = 1024
//
//	SpriteMagicHeader = "SP"
//
//	FileTypePAL FileType = iota
//	FileTypeRGBA
//)
//
//type SpriteFrame struct {
//	SpriteType FileType
//	Width      uintptr
//	Height     uintptr
//	Data       []byte
//}
//
//type SpriteFile struct {
//	Header             string
//	Version            float32
//	IndexedFrameCount  uint16
//	_IndexedFrameCount uint16
//	RGBAFrameCount     uint16
//	RGBAIndex          uint16
//	Palette            []byte
//	Frames             []*SpriteFrame
//}
//
//func Load(buf *bytes.Buffer) (file *SpriteFile, err error) {
//	file = &SpriteFile{
//		Palette: make([]byte, PaletteSize),
//	}
//
//	if err := file.parseHeader(buf); err != nil {
//		return nil, err
//	}
//
//	if file.Version < 2.1 {
//		if err = file.readIndexedFrames(buf); err != nil {
//			return nil, err
//		}
//	} else {
//		if err = file.readCompressedIndexedFrames(buf); err != nil {
//			return nil, err
//		}
//	}
//
//	if err = file.readRGBAFrames(buf); err != nil {
//		return nil, err
//	}
//
//	reader := bytes.NewReader(buf.Bytes())
//
//	if file.Version > 1.0 {
//		reader.Size()
//		currentPosition, err := reader.Seek(0, io.SeekCurrent)
//		if err != nil {
//			return nil, err
//		}
//
//		skip := int64((reader.Len() - 1024) - int(currentPosition))
//		if err = bytesUtil.SkipBytes(reader, skip); err != nil {
//			return nil, err
//		}
//
//		if _, err = io.ReadFull(io.LimitReader(reader, PaletteSize), file.Palette); err != nil {
//			return nil, err
//		}
//
//		log.Printf("len(file.Palette)=%v\n", len(file.Palette))
//	}
//
//	return file, nil
//}
//
//func (f *SpriteFile) parseHeader(buf io.Reader) error {
//	var header [2]byte
//	_ = binary.Read(buf, binary.LittleEndian, &header)
//
//	headerStr := string(header[:])
//	if headerStr != SpriteMagicHeader {
//		return fmt.Errorf("invalid header: %s\n", header)
//	}
//
//	var major, minor byte
//	_ = binary.Read(buf, binary.LittleEndian, &major)
//	_ = binary.Read(buf, binary.LittleEndian, &minor)
//
//	version, err := strconv.ParseFloat(fmt.Sprintf("%d.%d", major, minor), 32)
//	if err != nil {
//		return fmt.Errorf("invalid version: %s\n", strconv.FormatFloat(version, 'E', -1, 32))
//	}
//
//	var indexedFrameCount, rgbaFrameCount uint16
//	if err = binary.Read(buf, binary.LittleEndian, &indexedFrameCount); err != nil {
//		return err
//	}
//	f._IndexedFrameCount = indexedFrameCount + 0
//
//	if version > 1.1 {
//		if err = binary.Read(buf, binary.LittleEndian, &rgbaFrameCount); err != nil {
//			return err
//		}
//	}
//
//	f.Header = headerStr
//	f.Version = float32(version)
//	f.IndexedFrameCount = indexedFrameCount
//	f.Frames = make([]*SpriteFrame, indexedFrameCount+rgbaFrameCount)
//	f.RGBAIndex = f.IndexedFrameCount
//
//	return nil
//}
//
//// Parse .spr indexed images
//func (f *SpriteFile) readIndexedFrames(buf *bytes.Buffer) error {
//	for i := 0; i < int(f.IndexedFrameCount); i++ {
//		var (
//			width, height uint16
//			data          []byte
//		)
//
//		_ = binary.Read(buf, binary.LittleEndian, &width)
//		_ = binary.Read(buf, binary.LittleEndian, &height)
//
//		data = make([]byte, width*height)
//		reader := bufio.NewReader(bytes.NewReader(data))
//		_, err := io.ReadFull(io.LimitReader(reader, int64(width*height)), data)
//
//		if err != nil {
//			return errors.Wrap(err, "could not read indexed frames data")
//		}
//
//		f.Frames[i] = &SpriteFrame{
//			SpriteType: FileTypePAL,
//			Width:      uintptr(width),
//			Height:     uintptr(height),
//			Data:       data,
//		}
//	}
//
//	return nil
//}
//
//// Parse .spr indexed images encoded with RLE
//func (f *SpriteFile) readCompressedIndexedFrames(buf io.Reader) error {
//	for i := 0; i < int(f.IndexedFrameCount); i++ {
//		var (
//			width, height uint16
//			data          []byte
//		)
//
//		_ = binary.Read(buf, binary.LittleEndian, &width)
//		_ = binary.Read(buf, binary.LittleEndian, &height)
//
//		data, err := ioutil.ReadAll(io.LimitReader(buf, int64(width*height)))
//		if err != nil {
//			return errors.Wrap(err, "could not read indexed frames data")
//		}
//
//		f.Frames[i] = &SpriteFrame{
//			SpriteType: FileTypePAL,
//			Width:      uintptr(width),
//			Height:     uintptr(width),
//			Data:       data,
//		}
//	}
//
//	return nil
//}
//
//func (f *SpriteFile) readRGBAFrames(buf io.Reader) error {
//	for i := 0; i < int(f.RGBAFrameCount); i++ {
//		var (
//			width, height uint16
//			data          []byte
//		)
//
//		_ = binary.Read(buf, binary.LittleEndian, &width)
//		_ = binary.Read(buf, binary.LittleEndian, &height)
//
//		data, err := ioutil.ReadAll(io.LimitReader(buf, int64(width*height*4)))
//		if err != nil {
//			return errors.Wrap(err, "could not read indexed frames data")
//		}
//
//		f.Frames[i+int(f.RGBAIndex)] = &SpriteFrame{
//			SpriteType: FileTypeRGBA,
//			Width:      uintptr(width),
//			Height:     uintptr(width),
//			Data:       data,
//		}
//	}
//
//	return nil
//}
//
//func (f *SpriteFile) indexedToRGBA() error {
//	for i := 0; i < int(f.IndexedFrameCount); i++ {
//		frame := f.Frames[i]
//
//		if frame.SpriteType != FileTypePAL {
//			continue
//		}
//
//		for y := 0; y < int(frame.Height); y++ {
//			for x := 0; x < int(frame.Width); x++ {
//				width := int(frame.Width)
//				height := int(frame.Height)
//
//				idx1 := frame.Data[x+y*width] * 4
//				idx2 := (x + (height-y-1)*width) * 4
//
//				frame.Data[idx2+3] = f.Palette[idx1+0]
//				frame.Data[idx2+2] = f.Palette[idx1+1]
//				frame.Data[idx2+1] = f.Palette[idx1+2]
//
//				if idx1 != 0 {
//					frame.Data[idx2+0] = 255
//				}
//			}
//		}
//
//		frame.SpriteType = FileTypeRGBA
//	}
//
//	log.Printf("Palette len = %d\n", len(f.Palette))
//
//	return nil
//}
