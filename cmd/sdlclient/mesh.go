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

	numBuffers
)

type Mesh struct {
	vertices  []Vertex
	indices   []uint32
	vao, vbo  uint32
	vaBuffers [numBuffers]uint32
	drawCount int32
}

func NewMesh(vertices []Vertex) *Mesh {
	mesh := new(Mesh)
	mesh.vertices = vertices
	mesh.drawCount = int32(len(vertices))

	var positions []float32
	for _, v := range vertices {
		positions = append(positions, v.Position.X(), v.Position.Y(), v.Position.Z())
	}

	gl.GenBuffers(1, &mesh.vbo)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.BufferData(gl.ARRAY_BUFFER, 4*len(positions), gl.Ptr(positions), gl.STATIC_DRAW)

	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)
	gl.EnableVertexAttribArray(0)
	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vbo)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	return mesh
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.vao)
	gl.DrawArrays(gl.TRIANGLES, 0, m.drawCount)
	gl.BindVertexArray(0)
}
