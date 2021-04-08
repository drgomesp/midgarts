package main

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Vertex struct {
	Position  mgl32.Vec3
	Color     mgl32.Vec3
	TexCoords mgl32.Vec2
}

const (
	vbPosition = iota
	vbColor
	vbTexcoord
	vbIndex

	numBuffers
)

type Mesh struct {
	*Transform
	vertices  []Vertex
	indices   []uint32
	vao       uint32
	vaBuffers [numBuffers]uint32
	drawCount int32
}

func NewMesh(vertices []Vertex, indices []uint32) *Mesh {
	mesh := &Mesh{
		Transform: NewTransform(Origin),
		vertices:  vertices,
		indices:   indices,
		drawCount: int32(len(indices)),
	}

	gl.GenVertexArrays(1, &mesh.vao)
	gl.BindVertexArray(mesh.vao)

	var positions, colors, texCoords []float32
	for _, v := range vertices {
		positions = append(positions, v.Position.X(), v.Position.Y(), v.Position.Z())
		colors = append(colors, v.Color.X(), v.Color.Y(), v.Color.Z())
		texCoords = append(texCoords, v.TexCoords.X(), v.TexCoords.Y())
	}

	gl.GenBuffers(numBuffers, &mesh.vaBuffers[0])

	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vaBuffers[vbPosition])
	gl.BufferData(gl.ARRAY_BUFFER, len(positions)*4, gl.Ptr(positions), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(0)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, nil)

	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vaBuffers[vbColor])
	gl.BufferData(gl.ARRAY_BUFFER, len(colors)*4, gl.Ptr(colors), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(1)
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 0, nil)

	gl.BindBuffer(gl.ARRAY_BUFFER, mesh.vaBuffers[vbTexcoord])
	gl.BufferData(gl.ARRAY_BUFFER, len(texCoords)*4, gl.Ptr(texCoords), gl.STATIC_DRAW)

	gl.EnableVertexAttribArray(2)
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 0, nil)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, mesh.vaBuffers[vbIndex])
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(mesh.indices), gl.STATIC_DRAW)

	return mesh
}

func (m *Mesh) Draw() {
	gl.BindVertexArray(m.vao)
	gl.DrawElements(gl.TRIANGLES, m.drawCount, gl.UNSIGNED_INT, gl.Ptr(nil))
	gl.BindVertexArray(0)
}
