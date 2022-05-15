package graphic

import (
	"github.com/project-midgard/midgarts/internal/opengl"
)

type Geometry struct {
	gls                 *opengl.State
	handleVAO           uint32
	vbos                []*opengl.VBO
	indices             []uint32
	handleIndices       uint32
	shouldUpdateIndices bool
}

func NewGeometry() *Geometry {
	geometry := &Geometry{
		vbos:                nil,
		shouldUpdateIndices: true,
	}

	return geometry
}

func (g *Geometry) AddVBO(vbo *opengl.VBO) *Geometry {
	g.vbos = append(g.vbos, vbo)
	return g
}

func (g *Geometry) VBOs() []*opengl.VBO {
	return g.vbos
}

func (g *Geometry) Indices() []uint32 {
	return g.indices
}

func (g *Geometry) SetIndices(indices ...uint32) {
	g.indices = indices
}
