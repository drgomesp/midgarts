package opengl

import (
	"log"

	"github.com/go-gl/gl/v4.6-core/gl"
)

type Program struct {
	id uint32
}

func NewProgram() *Program {
	return &Program{}
}

func (p *Program) GetAttribLocation(name string) int32 {
	if p.id == 0 {
		log.Fatalf("program '%v' not loaded", p.id)
	}

	return gl.GetAttribLocation(p.id, gl.Str(name+"\x00"))
}

func (p *Program) ID() uint32 {
	return p.id
}
