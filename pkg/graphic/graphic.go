package graphic

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
)

const (
	OnePixelSize = 1.0 / 35.0
)

type Graphic struct {
	*Transform
	geometry   *Geometry
	renderMode uint32
}

func NewGraphic(geom *Geometry, renderMode uint32) *Graphic {
	return &Graphic{Transform: NewTransform(Origin), geometry: geom, renderMode: renderMode}
}

func (g *Graphic) Render(gls *opengl.State) {
	geom := g.geometry

	if geom.gls == nil {
		gl.GenVertexArrays(1, &geom.handleVAO)
		gl.GenBuffers(1, &geom.handleIndices)
		geom.gls = gls
	}

	gl.BindVertexArray(geom.handleVAO)
	for _, vbo := range geom.vbos {
		vbo.Load(gls)
	}

	if geom.shouldUpdateIndices {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, geom.handleIndices)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(geom.indices)*4, gl.Ptr(geom.indices), gl.STATIC_DRAW)
		geom.shouldUpdateIndices = false
	}

	gl.DrawElements(g.renderMode, int32(len(geom.Indices())*4), gl.UNSIGNED_INT, gl.Ptr(nil))
	gl.BindVertexArray(0)
}
