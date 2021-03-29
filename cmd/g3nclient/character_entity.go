package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
	"github.com/g3n/engine/graphic"
)

type CharacterEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*RenderComponent
}

func NewCharacterEntity(bodySprite *graphic.Sprite) *CharacterEntity {
	b := ecs.NewBasic()

	return &CharacterEntity{
		BasicEntity:     b.GetBasicEntity(),
		RenderComponent: &RenderComponent{Graphic: bodySprite},
	}
}

func (c *CharacterEntity) GetRenderComponent() *RenderComponent {
	return c.RenderComponent
}

func (c *CharacterEntity) GetSpaceComponent() *common.SpaceComponent {
	return c.SpaceComponent
}
