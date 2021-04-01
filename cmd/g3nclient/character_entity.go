package main

import (
	"log"

	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/internal/entity"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
)

type CharacterEntity struct {
	*ecs.BasicEntity
	*common.SpaceComponent
	*CharacterRenderComponent
	*CharacterAnimationComponent
	*CharacterControlComponent

	Animator      *Animator
	PlayMode      actionplaymode.Type
	Direction     directiontype.Type
	State         statetype.Type
	CurrentAction *entity.CharacterAction
}

func NewCharacterEntity(charSprite *CharacterSprite, animator *Animator) *CharacterEntity {
	b := ecs.NewBasic()

	e := &CharacterEntity{
		Animator:    animator,
		BasicEntity: b.GetBasicEntity(),
		CharacterRenderComponent: &CharacterRenderComponent{
			CharacterSprite: charSprite,
		},
		CharacterAnimationComponent: NewCharacterAnimationComponent(charSprite, animator, 0.1),
		CharacterControlComponent:   &CharacterControlComponent{},
		PlayMode:                    actionplaymode.Repeat,
		Direction:                   directiontype.South,
		State:                       statetype.Idle,
	}

	e.SetAction(statetype.Idle)
	return e
}

func (e *CharacterEntity) GetCharacterRenderComponent() *CharacterRenderComponent {
	return e.CharacterRenderComponent
}

func (e *CharacterEntity) GetCharacterAnimationComponent() *CharacterAnimationComponent {
	return e.CharacterAnimationComponent
}

func (e *CharacterEntity) GetCharacterControlComponent() *CharacterControlComponent {
	return e.CharacterControlComponent
}

func (e *CharacterEntity) GetSpaceComponent() *common.SpaceComponent {
	return e.SpaceComponent
}

func (e *CharacterEntity) GetCurrentAction() *entity.CharacterAction {
	return e.CurrentAction
}

func (e *CharacterEntity) GetDirection() directiontype.Type {
	return e.Direction
}

func (e *CharacterEntity) SetDirection(d directiontype.Type) {
	e.Direction = d
}

func (e *CharacterEntity) GetState() statetype.Type {
	return e.State
}

func (e *CharacterEntity) SetState(s statetype.Type) {
	e.State = s
}

func (e *CharacterEntity) GetPlayMode() actionplaymode.Type {
	return e.PlayMode
}

func (e *CharacterEntity) SetCurrentAction(action *entity.CharacterAction) {
	e.CurrentAction = action
}

func (e *CharacterEntity) SetAction(state statetype.Type) {
	e.State = state
	//e.CharacterAnimationComponent = &CharacterAnimationComponent{
	//	Animator:   e.Animator,
	//	Animations: map[string]*common.Animation{},
	//}
	e.CharacterAnimationComponent = NewCharacterAnimationComponent(
		e.GetCharacterRenderComponent().CharacterSprite,
		e.Animator,
		0.08,
	)
	e.CurrentAction = entity.NewCharacterAction(actionindex.GetActionIndex(e.State))
	anim := &common.Animation{Name: e.CurrentAction.Name, Frames: e.CurrentAction.Frames}
	e.CharacterAnimationComponent.AddAnimations([]*common.Animation{anim})
	e.CharacterAnimationComponent.AddDefaultAnimation(anim)
	e.CharacterAnimationComponent.CurrentAnimation = anim
	log.Printf("setting new action: state=%v\n", state)
}
