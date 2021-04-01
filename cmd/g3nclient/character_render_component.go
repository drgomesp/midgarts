package main

import (
	"github.com/EngoEngine/engo/common"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/window"
	"github.com/project-midgard/midgarts/internal/entity"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
)

type CharacterRenderFace interface {
	GetCharacterRenderComponent() *CharacterRenderComponent
	GetCharacterAnimationComponent() *CharacterAnimationComponent
	GetCharacterControlComponent() *CharacterControlComponent
}

type Character interface {
	common.BasicFace
	common.SpaceFace
	CharacterRenderFace

	GetPlayMode() actionplaymode.Type
	GetDirection() directiontype.Type
	SetDirection(directiontype.Type)
	GetState() statetype.Type
	SetState(statetype.Type)
	GetCurrentAction() *entity.CharacterAction
	SetCurrentAction(action *entity.CharacterAction)
	SetAction(statetype.Type)
}

type CharacterRenderComponent struct {
	CharacterSprite *CharacterSprite
	Scale           *math32.Vector3
}

type CharacterControlComponent struct {
	KeyState *window.KeyState
}

func (r *CharacterRenderComponent) GetCharacterRenderComponent() *CharacterRenderComponent {
	return r
}

func (ac *CharacterAnimationComponent) GetCharacterAnimationComponent() *CharacterAnimationComponent {
	return ac
}

func (ac *CharacterControlComponent) GetCharacterControlComponent() *CharacterControlComponent {
	return ac
}
