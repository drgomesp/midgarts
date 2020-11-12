package spr

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"log"
	"math"
	"strconv"

	bytesutil "github.com/project-midgard/midgarts/bytes"

	"github.com/pkg/errors"
)

type FileType int

const (
	HeaderSignature = "SP"
	PaletteSize     = 256 * 4

	SpriteFileTypePAL FileType = iota
	SpriteFileTypeRGBA
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
		Signature string
		Version   float32

		IndexedFrameCount uint16
		RGBAFrameCount    uint16
		RGBAIndex         uint16
	}

	Frames   []*SpriteFrame
	Palette  []byte
	Compiled bool
}

func Load(buf *bytes.Buffer) (f *SpriteFile, err error) {
	f = new(SpriteFile)
	reader := bytes.NewReader(buf.Bytes())

	if err := f.parseHeader(reader); err != nil {
		return nil, err
	}

	if f.Header.Version < 2.1 {
		return nil, fmt.Errorf("unsupported version %f\n", f.Header.Version)
	}

	if err = f.readCompressedIndexedFrames(reader); err != nil {
		return nil, err
	}

	if err = f.readRGBAFrames(reader); err != nil {
		return nil, err
	}

	f.parsePalette(buf)

	return f, nil
}

func (f *SpriteFile) parseHeader(buf io.ReadSeeker) error {
	log.Println("parsing header...")

	var signature [2]byte
	_ = binary.Read(buf, binary.LittleEndian, &signature)

	signatureStr := string(signature[:])
	if signatureStr != HeaderSignature {
		return fmt.Errorf("invalid signature: %s\n", signature)
	}

	var major, minor byte
	_ = binary.Read(buf, binary.LittleEndian, &minor)
	_ = binary.Read(buf, binary.LittleEndian, &major)

	version, err := strconv.ParseFloat(fmt.Sprintf("%d.%d", major, minor), 32)
	if err != nil {
		return errors.Wrapf(err, "invalid version: %s\n", strconv.FormatFloat(version, 'E', -1, 32))
	}

	var indexedFrameCount, rgbaFrameCount uint16
	_ = binary.Read(buf, binary.LittleEndian, &indexedFrameCount)

	if version > 1.1 {
		_ = binary.Read(buf, binary.LittleEndian, &rgbaFrameCount)
	}

	f.Header.Signature = signatureStr
	f.Header.Version = float32(version)
	f.Header.IndexedFrameCount = indexedFrameCount
	f.Header.RGBAFrameCount = rgbaFrameCount
	f.Header.RGBAIndex = indexedFrameCount
	f.Frames = make([]*SpriteFrame, indexedFrameCount+rgbaFrameCount)
	f.Palette = make([]byte, PaletteSize)

	return nil
}

// Parse .spr indexed images encoded with run-length encoding (RLE)
func (f *SpriteFile) readCompressedIndexedFrames(buf io.ReadSeeker) error {
	log.Println("reading compressed indexed frames...")

	for i := 0; i < int(f.Header.IndexedFrameCount); i++ {
		var (
			width, height, size, index, end uint16
			c, count                        byte
			data                            []byte
		)

		_ = binary.Read(buf, binary.LittleEndian, &width)
		_ = binary.Read(buf, binary.LittleEndian, &height)

		size = width * height
		data = make([]byte, size)
		index = 0

		_ = binary.Read(buf, binary.LittleEndian, &end)
		offset, _ := buf.Seek(0, io.SeekCurrent)
		end += uint16(offset)

		for {
			if uint16(offset) >= end || index >= size {
				break
			}

			offset, _ = buf.Seek(0, io.SeekCurrent)

			_ = binary.Read(buf, binary.LittleEndian, &c)
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
		}

		f.Frames[i] = &SpriteFrame{
			SpriteType: SpriteFileTypePAL,
			Width:      width,
			Height:     height,
			Data:       data,
		}
	}

	return nil
}

func (f *SpriteFile) readRGBAFrames(buf io.ReadSeeker) error {
	log.Println("reading rgba frames...")

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
		_ = binary.Read(buf, binary.LittleEndian, &data)

		f.Frames[i+int(f.Header.RGBAIndex)] = &SpriteFrame{
			SpriteType: SpriteFileTypeRGBA,
			Width:      width,
			Height:     height,
			Data:       data,
		}
	}

	return nil
}

func (f *SpriteFile) parsePalette(buf *bytes.Buffer) {
	reader := bytes.NewReader(buf.Bytes())
	_ = bytesutil.SkipBytes(reader, int64(reader.Len()-1024))
	_ = binary.Read(reader, binary.LittleEndian, &f.Palette)
}

func (f *SpriteFile) ToRGBA() {
	for i := 0; i < int(f.Header.IndexedFrameCount); i++ {
		frame := f.Frames[i]
		width, height := int(frame.Width), int(frame.Height)

		if frame.SpriteType != SpriteFileTypePAL {
			continue
		}

		buf := make([]byte, width*height*4)

		var idx1, idx2 int

		// reverse height
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx1 = int(frame.Data[x+y*width]) * 4
				idx2 = (x + (height-y-1)*width) * 4

				a := byte(0)

				if idx1 == 0 {
					a = 255
				}

				buf[idx2+3] = f.Palette[idx1+0]
				buf[idx2+2] = f.Palette[idx1+1]
				buf[idx2+1] = f.Palette[idx1+2]
				buf[idx2+0] = a
			}
		}

		frame.Data = buf
		frame.SpriteType = SpriteFileTypeRGBA
	}

	f.Header.IndexedFrameCount = 0
	f.Header.RGBAFrameCount = uint16(len(f.Frames))
	f.Header.RGBAIndex = 0
}

func (f *SpriteFile) Compile() {
	if !f.Compiled {
		for i := 0; i < len(f.Frames); i++ {
			f.Frames[i].Compile(f.Palette)
		}

		f.Compiled = true
	}
}

func (f *SpriteFrame) Compile(palette []byte) {
	if f.Compiled {
		return
	}

	var (
		glWidth, glHeight, startX, startY int
	)

	width, height := int(f.Width), int(f.Height)
	glWidth = int(math.Pow(2, math.Ceil(math.Log(float64(width))/math.Log(2))))
	glHeight = int(math.Pow(2, math.Ceil(math.Log(float64(height))/math.Log(2))))
	startX = int(math.Floor(float64(glWidth-width) * 0.5))
	startY = int(math.Floor(float64(glHeight-height) * 0.5))

	var buf []byte

	if f.SpriteType == SpriteFileTypePAL {
		buf = make([]byte, glWidth*glHeight)

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				index := ((y + startY) * glWidth) + (x + startX)
				buf[index] = f.Data[y*width+x]

				if palette[f.Data[y*width+x]*4] == 255 &&
					palette[f.Data[y*width+x]*4+2] == 255 &&
					palette[f.Data[y*width+x]*4+1] == 0 {
					buf[index] = 0
				}
			}
		}
	} else {
		buf = make([]byte, glWidth*glHeight*4)

		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				index := ((y+startY)*glWidth + (x + startX)) * 4
				buf[index+0] = f.Data[((height-y-1)*width+x)*4+3]
				buf[index+1] = f.Data[((height-y-1)*width+x)*4+2]
				buf[index+2] = f.Data[((height-y-1)*width+x)*4+1]
				buf[index+3] = f.Data[((height-y-1)*width+x)*4+0]
			}
		}
	}

	f.Compiled = true
	f.Data = buf
	f.Width = uint16(glWidth)
	f.Height = uint16(glHeight)
}

func (f *SpriteFile) Image(index int, mirrored bool) *image.RGBA {
	var (
		frame  = f.Frames[index]
		width  = int(frame.Width)
		height = int(frame.Height)
	)

	img := image.NewRGBA(image.Rect(0, 0, width, height))

	if frame.SpriteType == SpriteFileTypePAL {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				i := frame.Data[x+y*width] * 4
				a := byte(0)

				if i != 0 {
					a = 255
				}

				img.SetRGBA(x, y, color.RGBA{
					R: f.Palette[i+0],
					G: f.Palette[i+1],
					B: f.Palette[i+2],
					A: a,
				})
			}
		}
	} else if frame.SpriteType == SpriteFileTypeRGBA {
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {
				idx1 := (x + y*500) * 4
				idx2 := (x + y*width) * 4

				var a byte
				if idx1 != 0 {
					a = 255
				}

				img.SetRGBA(x, y, color.RGBA{
					R: frame.Data[idx2+3],
					G: frame.Data[idx2+2],
					B: frame.Data[idx2+1],
					A: a,
				})
			}
		}
	}

	return img
}
