package opengl

import (
	"github.com/go-gl/gl/v3.3-core/gl"
)

type VBOAttributeType int

const (
	VertexPosition = VBOAttributeType(iota)
	VertexColor
	VertexTexCoord

	NumVertexBuffers
)

var attribTypeNameMap = map[VBOAttributeType]string{
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
	gls        *State
	buf        []float32
	usage      uint32
	attributes []VBOAttribute
	buffers    [NumVertexBuffers]uint32
}

type VBOAttribute struct {
	Type        VBOAttributeType
	Name        string
	NumElements int
	ElementType uint32
}

func NewVBO(buf []float32) *VBO {
	vbo := &VBO{
		buf:        buf,
		usage:      gl.STATIC_DRAW,
		attributes: make([]VBOAttribute, 0),
	}

	return vbo
}

func (vbo *VBO) AddAttribute(t VBOAttributeType) *VBO {
	vbo.attributes = append(vbo.attributes, VBOAttribute{
		Type:        t,
		Name:        attribTypeNameMap[t],
		NumElements: attributeTypeSizes[t],
		ElementType: gl.FLOAT,
	})

	return vbo
}

func (vbo *VBO) StrideByteSize() int {
	stride := int(0)

	for _, attrib := range vbo.attributes {
		stride += attributeTypeSizes[attrib.Type] * attrib.NumElements
	}

	return stride
}

func (vbo *VBO) Load(gls *State) {
	if vbo.gls == nil {
		gl.GenBuffers(int32(NumVertexBuffers), &vbo.buffers[0])
		gl.BindBuffer(gl.ARRAY_BUFFER, vbo.buffers[0])

		if len(vbo.attributes) == 0 {
			return
		}

		for loc, attrib := range vbo.attributes {
			//loc := gls.program.GetAttribLocation(attrib.Name)
			//if loc < 0 {
			//	log.Fatalf("VBO attribute '%s' not found", attrib.Name)
			//}

			gl.EnableVertexAttribArray(uint32(loc))
			gl.VertexAttribPointer(uint32(loc), int32(attrib.NumElements), attrib.ElementType, false, int32(vbo.StrideByteSize()), nil)

			var data []float32
			for _, v := range vbo.buf {
				data = append(data, v)
			}

			gl.BindBuffer(gl.ARRAY_BUFFER, vbo.buffers[loc])
			gl.BufferData(gl.ARRAY_BUFFER, len(vbo.buf)*attributeTypeSizes[attrib.Type], gl.Ptr(data), vbo.usage)
		}
	}

	vbo.gls = gls
}
