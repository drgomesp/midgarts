package graphic

import (
	"github.com/drgomesp/midgarts/internal/opengl"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Sprite struct {
	positions []float32
	*Graphic
	Texture *Texture
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
}

func NewSprite(width, height float32, texture *Texture) *Sprite {
	s := &Sprite{
		positions: make([]float32, 12),
		Texture:   texture,
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

	s.Graphic = NewGraphic(geom, gl.TRIANGLES)
	return s
}
