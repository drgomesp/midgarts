package gnd

import (
	"bytes"
	"encoding/binary"
	"io"
	"log"

	"github.com/project-midgard/midgarts/pkg/bytesutil"
	"golang.org/x/text/encoding/charmap"
)

type GroundFile struct {
	Version        float32
	Width, Height  uint32
	Zoom           float32
	Textures       []string
	TextureIndices []int64
}

type LightMapData struct {
	PerCell uint32
	Count   uint32
	Data    []byte
}

func Load(buf *bytes.Buffer) (f *GroundFile, err error) {
	f = new(GroundFile)
	reader := bytes.NewReader(buf.Bytes())

	var signature [4]byte
	_ = binary.Read(reader, binary.LittleEndian, &signature)

	var a, b uint8
	_ = binary.Read(reader, binary.LittleEndian, &a)
	_ = binary.Read(reader, binary.LittleEndian, &b)
	f.Version = float32(a) + float32(b)/10

	var w, h uint32
	_ = binary.Read(reader, binary.LittleEndian, &w)
	_ = binary.Read(reader, binary.LittleEndian, &h)

	var zoom float32
	_ = binary.Read(reader, binary.LittleEndian, &zoom)

	f.Width = w
	f.Height = h
	f.Zoom = zoom

	err = f.loadTextures(reader)
	if err != nil {
		log.Fatal(err)
	}

	err = f.loadLightMaps(reader)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", f)

	return f, nil
}

func (f *GroundFile) loadTextures(buf io.Reader) error {
	var textureCount, texturePathLength uint32
	_ = binary.Read(buf, binary.LittleEndian, &textureCount)
	_ = binary.Read(buf, binary.LittleEndian, &texturePathLength)

	var textures []string
	lookUpList := make([]int64, textureCount)

	for i := 0; i < int(textureCount); i++ {
		name, err := bytesutil.ReadString(buf, int(texturePathLength))
		if err != nil {
			log.Fatal(err)
		}

		pos := -1
		for k, n := range textures {
			if name == n {
				pos = k
				break
			}
		}

		if pos == -1 {
			var decodedName string
			if decodedName, err = charmap.Windows1252.NewDecoder().String(
				name,
			); err != nil {
				panic(err)
			}

			textures = append(textures, decodedName)
			pos = len(textures) - 1
		}

		lookUpList[i] = int64(pos)
	}

	f.Textures = textures
	f.TextureIndices = lookUpList

	return nil
}

func (f *GroundFile) loadLightMaps(buf io.Reader) error {
	var count uint32
	_ = binary.Read(buf, binary.LittleEndian, &count)

	var perCellX, perCellY uint32
	_ = binary.Read(buf, binary.LittleEndian, &perCellX)
	_ = binary.Read(buf, binary.LittleEndian, &perCellY)

	var sizeCell uint32
	_ = binary.Read(buf, binary.LittleEndian, &sizeCell)

	perCell := perCellX * perCellY * sizeCell

	_ = perCell

	return nil
}
