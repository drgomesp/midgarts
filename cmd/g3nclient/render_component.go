package main

import (
	"github.com/EngoEngine/engo/common"
	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/math32"
)

type RenderFace interface {
	GetRenderComponent() *RenderComponent
}

type Character interface {
	common.BasicFace
	common.SpaceFace
	RenderFace
}

type RenderComponent struct {
	Graphic graphic.IGraphic
	Scale   *math32.Vector3
}

func (r *RenderComponent) GetRenderComponent() *RenderComponent {
	return r
}
