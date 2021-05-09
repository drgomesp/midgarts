package spr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"

	"github.com/project-midgard/midgarts/pkg/bytesutil"
)

type FileType uint32

const (
	HeaderSignature = "SP"
	PaletteSize     = 1024

	FileTypePAL  FileType = 0
	FileTypeRGBA FileType = 1
)

type SpriteFrame struct {
	SpriteType FileType
	Width      uint16
	Height     uint16
	Data       []byte
	Compiled   bool
}

type SpriteFile struct {
	Header struct {
		Signature          string
		Version            float32
		PalettedFrameCount uint16
		RGBAFrameCount     uint16
		RGBAIndex          uint16
	}

	Frames  []*SpriteFrame
	Palette [PaletteSize]byte
}

func Load(data []byte) (f *SpriteFile, err error) {
	f = new(SpriteFile)
	reader := bytes.NewReader(data)

	if err := f.parseHeader(reader); err != nil {
		return nil, err
	}

	if f.Header.Version < 2.1 {
		if err = f.readPalettedFrames(reader); err != nil {
			return nil, err
		}
	} else {
		if err = f.readPalettedFramesRLE(reader); err != nil {
			return nil, err
		}
	}

	if err = f.readRGBAFrames(reader); err != nil {
		return nil, err
	}

	if f.Header.Version > 1.0 {
		reader = bytes.NewReader(data)

		if err = f.parsePalette(int64(len(data)-PaletteSize), reader); err != nil {
			return nil, err
		}
	}

	return f, nil
}

func (f *SpriteFile) parsePalette(skip int64, buf io.ReadSeeker) error {
	_ = bytesutil.SkipBytes(buf, skip)
	return binary.Read(buf, binary.LittleEndian, &f.Palette)
}

func (f *SpriteFile) parseHeader(buf io.ReadSeeker) error {
	var signature [2]byte
	_ = binary.Read(buf, binary.LittleEndian, &signature)

	signatureStr := string(signature[:])
	if signatureStr != HeaderSignature {
		return fmt.Errorf("invalid signature: %s\n", signature)
	}

	var major, minor byte
	_ = binary.Read(buf, binary.LittleEndian, &major)
	_ = binary.Read(buf, binary.LittleEndian, &minor)
	version := float32(major)/10 + float32(minor)

	var palettedFrameCount, rgbaFrameCount uint16
	_ = binary.Read(buf, binary.LittleEndian, &palettedFrameCount)

	if version > 1.1 {
		_ = binary.Read(buf, binary.LittleEndian, &rgbaFrameCount)
	}

	f.Header.Signature = signatureStr
	f.Header.Version = version
	f.Header.PalettedFrameCount = palettedFrameCount
	f.Header.RGBAFrameCount = rgbaFrameCount
	f.Header.RGBAIndex = palettedFrameCount
	f.Frames = make([]*SpriteFrame, palettedFrameCount+rgbaFrameCount)

	return nil
}

func (f *SpriteFile) readPalettedFrames(buf io.ReadSeeker) error {
	var (
		width, height uint16
		size          int
		data          []byte
		palCount      = int(f.Header.PalettedFrameCount)
	)

	for i := 0; i < palCount; i++ {
		_ = binary.Read(buf, binary.LittleEndian, &width)
		_ = binary.Read(buf, binary.LittleEndian, &height)

		size = int(width * height)
		data = make([]byte, size)
		_ = binary.Read(buf, binary.LittleEndian, &data)

		f.Frames[i] = &SpriteFrame{
			SpriteType: FileTypePAL,
			Width:      width,
			Height:     height,
			Data:       data,
		}
	}

	return nil
}

// Parse .spr indexed images encoded with run-length encoding (RLE)
func (f *SpriteFile) readPalettedFramesRLE(buf io.ReadSeeker) error {
	var (
		width, height    uint16
		c, count         byte
		end, index, size int
		data             []byte
	)

	for i := 0; i < int(f.Header.PalettedFrameCount); i++ {
		_ = binary.Read(buf, binary.LittleEndian, &width)
		_ = binary.Read(buf, binary.LittleEndian, &height)

		size = int(width * height)
		data = make([]byte, size)
		index = 0

		var tmp uint16
		if err := binary.Read(buf, binary.LittleEndian, &tmp); err != nil {
			return err
		}

		offset, _ := buf.Seek(0, io.SeekCurrent)
		end = int(tmp) + int(offset)

		for int(offset) < end {
			if err := binary.Read(buf, binary.LittleEndian, &c); err != nil {
				return err
			}

			data[index] = c
			index++

			if c == 0 {
				_ = binary.Read(buf, binary.LittleEndian, &count)

				if count == 0 {
					data[index] = count
					index++
				} else {
					for j := 1; j < int(count); j++ {
						data[index] = c
						index++
					}
				}
			}

			offset, _ = buf.Seek(0, io.SeekCurrent)
		}

		f.Frames[i] = &SpriteFrame{
			SpriteType: FileTypePAL,
			Width:      width,
			Height:     height,
			Data:       data,
		}
	}

	return nil
}

func (f *SpriteFile) readRGBAFrames(buf io.ReadSeeker) error {
	for i := 0; i < int(f.Header.RGBAFrameCount); i++ {
		var (
			width, height uint16
			size          int
			data          []byte
		)

		_ = binary.Read(buf, binary.LittleEndian, &width)
		_ = binary.Read(buf, binary.LittleEndian, &height)

		size = int(width*height) * 4
		data = make([]byte, size)
		if err := binary.Read(buf, binary.LittleEndian, &data); err != nil {
			return err
		}

		f.Frames[i+int(f.Header.RGBAIndex)] = &SpriteFrame{
			SpriteType: FileTypeRGBA,
			Width:      width,
			Height:     height,
			Data:       data,
		}
	}

	return nil
}

func (f *SpriteFile) ImageAt(index int) *image.RGBA {
	var (
		frame  = f.Frames[index]
		width  = int(frame.Width)
		height = int(frame.Height)
		data   = frame.Data
	)

	if width <= 0 || height <= 0 {
		return nil
	}

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	if frame.SpriteType == FileTypeRGBA {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				i := (x + y*width) * 4

				img.Set(x, y, color.RGBA{
					R: data[i+3],
					G: data[i+2],
					B: data[i+1],
					A: data[i+0],
				})
			}
		}
	} else {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				i := int(frame.Data[x+y*width]) * 4
				var a byte

				if i != 0 {
					a = 255
				}

				img.Set(x, y, color.RGBA{
					R: f.Palette[i+0],
					G: f.Palette[i+1],
					B: f.Palette[i+2],
					A: a,
				})
			}
		}
	}

	return img
}
