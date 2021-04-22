package graphic

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	opengl2 "github.com/project-midgard/midgarts/internal/opengl"
)

type Plane struct {
	*Graphic
}

func NewPlane(width, height float32) *Plane {
	geom := NewGeometry()
	w := width / 2
	h := height / 2

	positions := []float32{
		+w, +h, 0, // top-left
		-w, -h, 0, // bottom-right
		+w, -h, 0, // bottom-left
		-w, +h, 0, // top-right
	}
	colors := []float32{
		1, 0, 0,
		0, 1, 0,
		0, 0, 1,
		1, 1, 1,
	}
	texCoords := []float32{
		0, 0,
		1, 1,
		0, 1,
		1, 0,
	}

	geom.AddVBO(opengl2.NewVBO([opengl2.NumVertexAttributes][]float32{
		positions,
		colors,
		texCoords,
	}).AddAttribute(opengl2.VertexPosition).
		AddAttribute(opengl2.VertexColor).
		AddAttribute(opengl2.VertexTexCoord)).
		SetIndices(0, 1, 2, 3, 1, 0)

	return &Plane{
		Graphic: NewGraphic(geom, gl.TRIANGLES),
	}
}
