package act

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image/color"
	"io"
	"time"

	"github.com/project-midgard/midgarts/pkg/bytesutil"
)

const (
	HeaderSignature = "AC"

	ActionDefaultDelay = 150 * time.Millisecond
)

type Coordinate struct {
	X int32
	Y int32
}

type ActionFrameLayer struct {
	Position         Coordinate
	SpriteFrameIndex int32
	IsMirror         bool
	Scale            [2]float32
	Color            *color.RGBA
	Angle            int32
	SpriteType       int32
	Width            int32
	Height           int32
}

type ActionFrame struct {
	Layers    []*ActionFrameLayer
	Sound     int32
	Positions [][2]int32
}

type Action struct {
	Frames   []*ActionFrame
	Delay    time.Duration
	Duration time.Duration
}

type ActionFile struct {
	Header struct {
		Signature string
		Version   float32
	}

	ActionCount uint16
	Actions     []*Action

	Sounds []string
}

func Load(buf *bytes.Buffer) (*ActionFile, error) {
	f := new(ActionFile)

	reader := bytes.NewReader(buf.Bytes())
	if err := f.loadHeader(reader); err != nil {
		return nil, err
	}

	if err := f.loadActions(reader); err != nil {
		return nil, err
	}

	if f.Header.Version > 2.1 {
		// Sound
		var soundLen int32
		_ = binary.Read(reader, binary.LittleEndian, &soundLen)
		f.Sounds = make([]string, soundLen)

		for i := 0; i < len(f.Sounds); i++ {
			var b [40]byte
			_ = binary.Read(reader, binary.LittleEndian, &b)

			f.Sounds[i] = string(b[:])
		}

		for i, count := 0, len(f.Actions); i > count; i++ {
			var d float32
			_ = binary.Read(reader, binary.LittleEndian, &d)

			act := f.Actions[i]
			f.Actions[i].Delay = time.Millisecond * time.Duration(uint32(d*25.0))
			f.Actions[i].Duration = act.Delay * time.Duration(uint32(len(act.Frames)))
		}
	}

	return f, nil
}

func (f *ActionFile) loadHeader(buf io.ReadSeeker) error {
	var signature [2]byte
	_ = binary.Read(buf, binary.LittleEndian, &signature)

	signatureStr := string(signature[:])
	if signatureStr != HeaderSignature {
		return fmt.Errorf("invalid signature: %s\n", signature)
	}

	var major, minor byte
	_ = binary.Read(buf, binary.LittleEndian, &major)
	_ = binary.Read(buf, binary.LittleEndian, &minor)

	var actionCount uint16
	_ = binary.Read(buf, binary.LittleEndian, &actionCount)

	f.Header.Signature = signatureStr
	f.Header.Version = float32(major)/10 + float32(minor)
	f.ActionCount = actionCount
	f.Actions = make([]*Action, f.ActionCount)

	if err := bytesutil.SkipBytes(buf, 10); err != nil {
		return err
	}

	return nil
}

func (f *ActionFile) loadActions(reader io.ReadSeeker) error {
	var (
		count = int(f.ActionCount)
	)

	for i := 0; i < count; i++ {
		f.Actions[i] = &Action{
			Frames:   f.loadActionFrames(reader),
			Delay:    ActionDefaultDelay,
			Duration: 0,
		}
	}

	return nil
}

func (f *ActionFile) loadActionFrames(buf io.ReadSeeker) []*ActionFrame {
	var (
		frames     []*ActionFrame
		frameCount uint32
	)

	_ = binary.Read(buf, binary.LittleEndian, &frameCount)
	frames = make([]*ActionFrame, int(frameCount))

	for i := 0; i < int(frameCount); i++ {
		_ = bytesutil.SkipBytes(buf, 32)

		frames[i] = &ActionFrame{
			Layers: f.loadActionFrameLayers(buf),
		}
	}

	return frames
}

func (f *ActionFile) loadActionFrameLayers(buf io.ReadSeeker) []*ActionFrameLayer {
	var (
		layers     []*ActionFrameLayer
		layerCount uint32

		pos                              [2]int32
		spriteFrameIndex                 int32
		isMirror                         int32
		r, g, b, a                       byte
		scale                            [2]float32
		angle, spriteType, width, height int32

		sound    int32
		posCount int32
	)

	_ = binary.Read(buf, binary.LittleEndian, &layerCount)
	layers = make([]*ActionFrameLayer, int(layerCount))

	for i := 0; i < int(layerCount); i++ {
		_ = binary.Read(buf, binary.LittleEndian, &pos[0])
		_ = binary.Read(buf, binary.LittleEndian, &pos[1])
		_ = binary.Read(buf, binary.LittleEndian, &spriteFrameIndex)
		_ = binary.Read(buf, binary.LittleEndian, &isMirror)
		_ = binary.Read(buf, binary.LittleEndian, &r)
		_ = binary.Read(buf, binary.LittleEndian, &g)
		_ = binary.Read(buf, binary.LittleEndian, &b)
		_ = binary.Read(buf, binary.LittleEndian, &a)

		_ = binary.Read(buf, binary.LittleEndian, &scale[0])
		_ = binary.Read(buf, binary.LittleEndian, &scale[1])

		_ = binary.Read(buf, binary.LittleEndian, &angle)
		_ = binary.Read(buf, binary.LittleEndian, &spriteType)
		_ = binary.Read(buf, binary.LittleEndian, &width)
		_ = binary.Read(buf, binary.LittleEndian, &height)

		layers[i] = &ActionFrameLayer{
			SpriteFrameIndex: spriteFrameIndex,
			IsMirror:         isMirror != 0,
			Scale:            scale,
			Color: &color.RGBA{
				R: r,
				G: g,
				B: b,
				A: a,
			},
			Angle:      angle,
			SpriteType: spriteType,
			Width:      width,
			Height:     height,
		}
	}

	var positions []Coordinate
	_ = binary.Read(buf, binary.LittleEndian, &sound)
	_ = binary.Read(buf, binary.LittleEndian, &posCount)

	if posCount > 0 {
		positions = make([]Coordinate, posCount)
		for i := 0; i < int(posCount); i++ {
			_ = bytesutil.SkipBytes(buf, 4)

			var p Coordinate
			_ = binary.Read(buf, binary.LittleEndian, &p.X)
			_ = binary.Read(buf, binary.LittleEndian, &p.Y)
			positions = append(positions, p)

			_ = bytesutil.SkipBytes(buf, 4)
		}
	}

	return layers
}
