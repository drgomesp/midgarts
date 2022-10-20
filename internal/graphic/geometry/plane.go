package geometry

import (
	"github.com/go-gl/gl/v4.6-core/gl"

	"github.com/project-midgard/midgarts/internal/graphic"
	"github.com/project-midgard/midgarts/internal/opengl"
)

const (
	OnePixelSize = 1.0 / 35.0
)

type Plane struct {
	*graphic.Graphic

	Geometry *graphic.Geometry
	Texture  *graphic.Texture

	Width, Height float32

	colors    []float32
	positions []float32
}

func (s *Plane) SetBounds(width, height float32) {
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

func (s *Plane) SetColors(colors []float32) {
	s.colors = colors
}

func (s *Plane) SetTexture(text *graphic.Texture) {
	s.Texture = text
	s.Texture.Bind(0)
}

func NewPlane(width, height float32, texture *graphic.Texture) *Plane {
	plane := &Plane{
		Texture:   texture,
		Width:     width,
		Height:    height,
		positions: make([]float32, 12),
	}

	plane.SetBounds(width, height)
	if plane.Texture != nil {
		plane.SetColors([]float32{
			1, 1, 1,
			1, 1, 1,
			1, 1, 1,
			1, 1, 1,
		})
	} else {
		plane.SetColors([]float32{
			1, 0, 1,
			1, 0, 1,
			1, 0, 1,
			1, 0, 1,
		})
	}

	texCoords := []float32{
		0, 0,
		1, 1,
		0, 1,
		1, 0,
	}

	geom := graphic.NewGeometry()
	geom.AddVBO(opengl.NewVBO([opengl.NumVertexAttributes][]float32{
		plane.positions,
		plane.colors,
		texCoords,
	}).AddAttribute(opengl.VertexPosition).
		AddAttribute(opengl.VertexColor).
		AddAttribute(opengl.VertexTexCoord),
	).SetIndices(0, 1, 2, 3, 1, 0)

	plane.Geometry = geom
	plane.Graphic = graphic.NewGraphic(geom, gl.TRIANGLES)

	return plane
}
