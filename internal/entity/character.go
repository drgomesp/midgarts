package entity

import (
	"github.com/EngoEngine/ecs"
	"github.com/project-midgard/midgarts/internal/component"
	"github.com/project-midgard/midgarts/pkg/character"
	"github.com/project-midgard/midgarts/pkg/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/character/statetype"
	"github.com/project-midgard/midgarts/pkg/graphic"
)

type Character struct {
	*graphic.Transform
	*ecs.BasicEntity
	*component.CharacterAttachmentComponent
	*component.CharacterStateComponent
	*component.CharacterSpriteRenderInfoComponent

	HeadIndex     int
	Gender        character.GenderType
	JobSpriteID   jobspriteid.Type
	IsMounted     bool
	MovementSpeed float64
}

func NewCharacter(gender character.GenderType, jobSpriteID jobspriteid.Type, headIndex int) *Character {
	b := ecs.NewBasic()
	c := &Character{
		BasicEntity: &b,
		CharacterStateComponent: &component.CharacterStateComponent{
			PlayMode: actionplaymode.Repeat,
			State:    statetype.Idle,
		},
		CharacterSpriteRenderInfoComponent: component.NewCharacterSpriteRenderInfoComponent(),
		Transform:                          graphic.NewTransform(graphic.Origin),
		Gender:                             gender,
		JobSpriteID:                        jobSpriteID,
		HeadIndex:                          headIndex,
		IsMounted:                          true,
		MovementSpeed:                      1.25,
	}

	return c
}

func (c *Character) SetCharacterStateComponent(component *component.CharacterStateComponent) {
	c.CharacterStateComponent = component
}

func (c *Character) SetCharacterAttachmentComponent(component *component.CharacterAttachmentComponent) {
	c.CharacterAttachmentComponent = component
}

func (c *Character) GetCharacterStateComponent() *component.CharacterStateComponent {
	return c.CharacterStateComponent
}

func (c *Character) GetCharacterAttachmentComponent() *component.CharacterAttachmentComponent {
	return c.CharacterAttachmentComponent
}

func (c *Character) GetCharacterSpriteRenderInfoComponent() *component.CharacterSpriteRenderInfoComponent {
	return c.CharacterSpriteRenderInfoComponent
}

func (c *Character) SetState(state statetype.Type) {
	c.PreviousState = c.State
	c.State = state
}
