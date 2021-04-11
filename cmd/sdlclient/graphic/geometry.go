package graphic

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/project-midgard/midgarts/cmd/sdlclient/array"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
)

type Geometry struct {
	gls                 *opengl.State
	handleVAO           uint32
	vbos                []*opengl.VBO
	indices             array.Uint32
	handleIndices       uint32
	shouldUpdateIndices bool
}

func NewGeometry() *Geometry {
	geometry := &Geometry{
		vbos:                make([]*opengl.VBO, 0),
		indices:             make(array.Uint32, 0),
		shouldUpdateIndices: true,
	}

	return geometry
}

func (g *Geometry) AddVBO(vbo *opengl.VBO) {
	g.vbos = append(g.vbos, vbo)
}

func (g *Geometry) Indices() array.Uint32 {
	return g.indices
}

func (g *Geometry) SetIndices(indices array.Uint32) {
	g.indices = indices
}

func (g *Geometry) PreRender(gls *opengl.State, _ *Camera) {
	if g.gls == nil {
		gl.GenVertexArrays(1, &g.handleVAO)
		gl.GenBuffers(1, &g.handleIndices)
		g.gls = gls
	}

	gl.BindVertexArray(g.handleVAO)
	for _, vbo := range g.vbos {
		vbo.Load(gls)
	}

	if g.shouldUpdateIndices {
		gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, uint32(opengl.NumVertexBuffers))
		gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, g.indices.Bytes(), gl.Ptr(g.indices.ToUint32()), gl.STATIC_DRAW)
		g.shouldUpdateIndices = false
	}
}
