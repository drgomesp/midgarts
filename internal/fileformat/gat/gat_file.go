package gat

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/project-midgard/midgarts/internal/romap"
)

const HeaderSignature = "GRAT"

const (
	None     = romap.CellTypeNone
	Walkable = romap.CellTypeWalkable
	Water    = romap.CellTypeWater
	Snipable = romap.CellTypeSnipable
)

var TypeTable = [7]romap.CellType{
	Walkable | Snipable,         // walkable ground
	None,                        // non-walkable ground
	Walkable | Snipable,         // ???
	Walkable | Snipable | Water, // walkable water
	Walkable | Snipable,         // ???
	Snipable,                    // gat (snipable)
	Walkable | Snipable,         // ???
}

type Cell struct {
	Cells    [4]float32
	CellType romap.CellType
}

type GroundAltitudeFile struct {
	Version float32
	Width   uint32
	Height  uint32
	Cells   []Cell
}

type BlockingRectangle struct {
}

func Load(data []byte) (f *GroundAltitudeFile, err error) {
	f = new(GroundAltitudeFile)
	reader := bytes.NewReader(data)

	var signature [4]byte
	_ = binary.Read(reader, binary.LittleEndian, &signature)

	if string(signature[:]) != HeaderSignature {
		return nil, fmt.Errorf("invalid file header signature: %s", signature)
	}

	var a, b byte
	_ = binary.Read(reader, binary.LittleEndian, &a)
	_ = binary.Read(reader, binary.LittleEndian, &b)
	f.Version = float32(a) + float32(b)/10

	var w, h uint32
	_ = binary.Read(reader, binary.LittleEndian, &w)
	_ = binary.Read(reader, binary.LittleEndian, &h)
	f.Width = w
	f.Height = h

	cellCount := w * h
	f.Cells = make([]Cell, cellCount)
	for i := 0; i < int(cellCount); i++ {
		var h1, h2, h3, h4 float32
		_ = binary.Read(reader, binary.LittleEndian, &h1)
		_ = binary.Read(reader, binary.LittleEndian, &h2)
		_ = binary.Read(reader, binary.LittleEndian, &h3)
		_ = binary.Read(reader, binary.LittleEndian, &h4)

		var t uint32
		_ = binary.Read(reader, binary.LittleEndian, &t)

		f.Cells[i] = Cell{
			Cells:    [4]float32{h1, h2, h3, h4},
			CellType: TypeTable[t],
		}
	}

	//log.Printf("gat.GroundAltitudeFile(%v)", f)

	return nil, nil
}
