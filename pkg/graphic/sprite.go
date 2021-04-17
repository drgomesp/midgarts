package graphic

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
)

type Sprite struct {
	*Graphic
	Texture *Texture
}

func NewSprite(width, height float32, texture *Texture) *Sprite {
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
		positions,
		colors,
		texCoords,
	}).AddAttribute(opengl.VertexPosition).
		AddAttribute(opengl.VertexColor).
		AddAttribute(opengl.VertexTexCoord),
	).SetIndices(0, 1, 2, 3, 1, 0)

	return &Sprite{
		Graphic: NewGraphic(geom, gl.TRIANGLES),
		Texture: texture,
	}
}
