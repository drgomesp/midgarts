package opengl

import (
	"unsafe"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type VBOAttributeType int

const (
	VertexPosition = VBOAttributeType(iota)
	VertexColor
	VertexTexCoord

	NumVertexAttributes
)

var attribTypeNames = map[VBOAttributeType]string{
	VertexPosition: "VertexPosition",
	VertexColor:    "VertexColor",
	VertexTexCoord: "VertexTexCoord",
}

var attributeTypeSizes = map[VBOAttributeType]int{
	VertexPosition: 3,
	VertexColor:    3,
	VertexTexCoord: 2,
}

type VBO struct {
	gls           *State
	initialized   bool
	buffers       [NumVertexAttributes][]float32
	usage         uint32
	attributes    []VBOAttribute
	handles       [NumVertexAttributes]uint32
	shouldUpdate  bool
	buffAllocated bool
}

type VBOAttribute struct {
	Type        VBOAttributeType
	Name        string
	NumElements int32
	ByteOffset  uint32
	ElementType uint32
}

func NewVBO(buffers [NumVertexAttributes][]float32) *VBO {
	vbo := &VBO{
		gls:          nil,
		buffers:      buffers,
		usage:        gl.STATIC_DRAW,
		attributes:   make([]VBOAttribute, 0),
		shouldUpdate: true,
	}

	return vbo
}

func (vbo *VBO) AddAttribute(t VBOAttributeType) *VBO {
	vbo.attributes = append(vbo.attributes, VBOAttribute{
		Type:        t,
		Name:        attribTypeNames[t],
		NumElements: int32(attributeTypeSizes[t]),
		ByteOffset:  0,
		ElementType: gl.FLOAT,
	})

	return vbo
}

func (vbo *VBO) Load(gls *State) {
	if len(vbo.attributes) == 0 {
		return
	}

	if !vbo.buffAllocated {
		for loc, _ := range vbo.attributes {
			gl.GenBuffers(int32(NumVertexAttributes), &vbo.handles[loc])
		}
		vbo.buffAllocated = true
	}

	for loc, attrib := range vbo.attributes {
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo.handles[loc])
		gl.BufferData(gl.ARRAY_BUFFER, len(vbo.buffers[loc])*4, gl.Ptr(vbo.buffers[loc]), vbo.usage)

		gl.EnableVertexAttribArray(uint32(loc))
		gl.VertexAttribPointer(
			uint32(loc),
			attrib.NumElements,
			attrib.ElementType,
			false,
			0,
			unsafe.Pointer(uintptr(attrib.ByteOffset)),
		)
	}

	vbo.gls = gls
	if !vbo.shouldUpdate {
		return
	}

	vbo.shouldUpdate = false
}
