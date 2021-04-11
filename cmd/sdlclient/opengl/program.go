package opengl

import "github.com/go-gl/gl/v3.3-core/gl"

type Program struct {
	id uint32
}

func NewProgram() *Program {
	return &Program{}
}

func (p *Program) GetAttribLocation(name string) int32 {
	return gl.GetAttribLocation(p.id, gl.Str(name+"\x00"))
}

func (p *Program) ID() uint32 {
	return p.id
}
