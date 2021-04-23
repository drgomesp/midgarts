package entity

import (
	"github.com/EngoEngine/ecs"
	"github.com/project-midgard/midgarts/internal/component"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
	"github.com/project-midgard/midgarts/pkg/graphic"
)

type Character struct {
	*ecs.BasicEntity
	*component.CharacterActionComponent
	*component.CharacterAttachmentComponent
	*graphic.Transform

	HeadIndex   int
	Gender      character.GenderType
	JobSpriteID jobspriteid.Type
	IsMounted   bool
}

func NewCharacter(gender character.GenderType, jobSpriteID jobspriteid.Type, headIndex int) *Character {
	b := ecs.NewBasic()
	return &Character{
		BasicEntity: &b,
		CharacterActionComponent: &component.CharacterActionComponent{
			PlayMode:  actionplaymode.Repeat,
			Action:    actionindex.Idle,
			Direction: directiontype.South,
			State:     statetype.Idle,
		},
		Transform:   graphic.NewTransform(graphic.Origin),
		Gender:      gender,
		JobSpriteID: jobSpriteID,
		HeadIndex:   headIndex,
		IsMounted:   true,
	}
}

func (c *Character) SetCharacterActionComponent(component *component.CharacterActionComponent) {
	c.CharacterActionComponent = component
}

func (c *Character) SetCharacterAttachmentComponent(component *component.CharacterAttachmentComponent) {
	c.CharacterAttachmentComponent = component
}

func (c *Character) GetCharacterActionComponent() *component.CharacterActionComponent {
	return c.CharacterActionComponent
}

func (c *Character) GetCharacterAttachmentComponent() *component.CharacterAttachmentComponent {
	return c.CharacterAttachmentComponent
}
