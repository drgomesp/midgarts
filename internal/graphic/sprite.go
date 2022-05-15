package graphic

import (
	"github.com/go-gl/gl/v4.6-core/gl"

	"github.com/project-midgard/midgarts/internal/opengl"
)

const (
	OnePixelSize = 1.0 / 35.0
)

type Sprite struct {
	*Graphic

	Geometry *Geometry
	Texture  *Texture

	Width, Height float32
	positions     []float32
}

func (s *Sprite) SetBounds(width, height float32) {
	w := width / 2
	h := height / 2

	s.positions[0] = -w
	s.positions[3] = w
	s.positions[6] = -w
	s.positions[9] = w

	s.positions[1] = h
	s.positions[4] = -h
	s.positions[7] = -h
	s.positions[10] = h
}

func (s *Sprite) SetTexture(text *Texture) {
	s.Texture = text
	s.Texture.Bind(0)
}

func NewSprite(width, height float32, texture *Texture) *Sprite {
	s := &Sprite{
		Texture:   texture,
		Width:     width,
		Height:    height,
		positions: make([]float32, 12),
	}

	s.SetBounds(width, height)

	geom := NewGeometry()

	colors := []float32{
		1, 1, 1,
		1, 1, 1,
		1, 1, 1,
		1, 1, 1,
	}
	texCoords := []float32{
		0, 0,
		1, 1,
		0, 1,
		1, 0,
	}

	geom.AddVBO(opengl.NewVBO([opengl.NumVertexAttributes][]float32{
		s.positions,
		colors,
		texCoords,
	}).AddAttribute(opengl.VertexPosition).
		AddAttribute(opengl.VertexColor).
		AddAttribute(opengl.VertexTexCoord),
	).SetIndices(0, 1, 2, 3, 1, 0)

	s.Geometry = geom
	s.Graphic = NewGraphic(geom, gl.TRIANGLES)

	return s
}
