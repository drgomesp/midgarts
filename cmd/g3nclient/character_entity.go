package main

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
)

type CharacterEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*RenderComponent
}

func NewCharacterEntity(charSprite *CharacterSprite) *CharacterEntity {
	b := ecs.NewBasic()

	return &CharacterEntity{
		BasicEntity:     b.GetBasicEntity(),
		RenderComponent: &RenderComponent{Graphic: charSprite.bodySprite},
	}
}

func (c *CharacterEntity) GetRenderComponent() *RenderComponent {
	return c.RenderComponent
}

func (c *CharacterEntity) GetSpaceComponent() *common.SpaceComponent {
	return c.SpaceComponent
}
