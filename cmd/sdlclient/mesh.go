package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Vertex struct {
	Position mgl32.Vec3
}

const (
	vbPosition = iota
	vbIndex
	numBuffers
)

type Mesh struct {
	vertices  []Vertex
	indices   []uint32
	vao       uint32
	vaBuffers [numBuffers]uint32
	drawCount int32
}

func NewMesh(vertices []Vertex, indices []uint32) *Mesh {
	mesh := &Mesh{
		vertices:  vertices,
		indices:   indices,
		drawCount: int32(len(indices)),
	}

	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	var positions []float32
	for _, v := range vertices {
		positions = append(positions, v.Position.X(), v.Position.Y(), v.Position.Z())
	}

	gl.GenBuffers(numBuffers, &mesh.vaBuffers[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vaBuffers[vbPosition])
	gl.BufferData(gl.ARRAY_BUFFER, len(positions)*4, gl.Ptr(positions), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.vaBuffers[vbIndex])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(mesh.indices), gl.STATIC_DRAW)

	return mesh
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, m.drawCount, gl.UNSIGNED_INT, gl.Ptr(nil))
	gl.BindVertexArray(0)
}
