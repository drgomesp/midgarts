package graphic

import (
	"github.com/go-gl/gl/v4.6-core/gl"

	"github.com/project-midgard/midgarts/internal/opengl"
)

type Plane struct {
	*Graphic

	Geometry *Geometry

	Width, Height float32
	positions     []float32
}

func (p *Plane) SetBounds(width, height float32) {
	w := width / 2
	h := height / 2

	p.positions[0] = -w
	p.positions[3] = w
	p.positions[6] = -w
	p.positions[9] = w

	p.positions[1] = h
	p.positions[4] = -h
	p.positions[7] = -h
	p.positions[10] = h
}

func NewPlane(width, height float32, texture *Texture) *Plane {
	p := &Plane{
		Width:     width,
		Height:    height,
		positions: make([]float32, 12),
	}

	p.SetBounds(width, height)
	geom := NewGeometry()

	colors := []float32{
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
		1, 0, 1,
	}
	texCoords := []float32{
		0, 0,
		1, 1,
		0, 1,
		1, 0,
	}

	geom.AddVBO(opengl.NewVBO([opengl.NumVertexAttributes][]float32{
		p.positions,
		colors,
		texCoords,
	}).AddAttribute(opengl.VertexPosition).
		AddAttribute(opengl.VertexColor).
		AddAttribute(opengl.VertexTexCoord),
	).SetIndices(0, 1, 2, 3, 1, 0)

	p.Geometry = geom
	p.Graphic = NewGraphic(geom, gl.TRIANGLES)

	return p
}
