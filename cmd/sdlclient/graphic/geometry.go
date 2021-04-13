package graphic

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/project-midgard/midgarts/cmd/sdlclient/array"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
)

type Geometry struct {
	*Transform
	gls                 *opengl.State
	handleVAO           uint32
	vbos                []*opengl.VBO
	indices             []uint32
	handleIndices       uint32
	shouldUpdateIndices bool
	renderMode          uint32
}

func NewGeometry(renderMode uint32) *Geometry {
	geometry := &Geometry{
		Transform:           NewTransform(Origin),
		vbos:                nil,
		shouldUpdateIndices: true,
		renderMode:          renderMode,
	}

	return geometry
}

func (g *Geometry) AddVBO(vbo *opengl.VBO) *Geometry {
	g.vbos = append(g.vbos, vbo)
	return g
}

func (g *Geometry) Indices() array.Uint32 {
	return g.indices
}

func (g *Geometry) SetIndices(indices []uint32) {
	g.indices = indices
}

func (g *Geometry) Render(gls *opengl.State, _ *Camera) {
	if g.gls == nil {
		gl.GenVertexArrays(1, &g.handleVAO)
		gl.GenBuffers(1, &g.handleIndices)
		g.gls = gls
	}

	gl.BindVertexArray(g.handleVAO)
	for _, vbo := range g.vbos {
		vbo.Load(g.gls)
	}

	if g.shouldUpdateIndices {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, g.handleIndices)
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(g.indices)*4, gl.Ptr(g.indices), gl.STATIC_DRAW)
		g.shouldUpdateIndices = false
	}

	gl.DrawElements(g.renderMode, g.Indices().Size(), gl.UNSIGNED_INT, gl.Ptr(nil))
	gl.BindVertexArray(0)
}
