package client

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/pkg/client/component"
	"github.com/project-midgard/midgarts/pkg/client/graphics"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/jobid"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	uuid "github.com/satori/go.uuid"
)

type CharacterEntity struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	component.CharacterAnimationComponent

	id string

	// spritesheetResource should come from a shared system later, but for now its ok
	// for every character to hold its own spritesheetResource.
	SpritesheetResource *graphics.SpritesheetResource
	ActionFile          *act.ActionFile
	CurrentAction       *CharacterAction
	PlayMode            actionplaymode.Type

	Gender    character.GenderType
	Job       jobid.Type
	State     statetype.Type
	Direction directiontype.Type
}

func NewCharacterEntity(spritesheetResource *graphics.SpritesheetResource, actFile *act.ActionFile, gender character.GenderType, job jobid.Type) *CharacterEntity {
	return &CharacterEntity{
		BasicEntity:         ecs.NewBasic(),
		id:                  uuid.NewV4().String(),
		SpritesheetResource: spritesheetResource,
		ActionFile:          actFile,
		Gender:              gender,
		Job:                 job,
		State:               statetype.Idle,
		Direction:           directiontype.South,
		PlayMode:            actionplaymode.Repeat,
	}
}

func (c *CharacterEntity) SetAction(state statetype.Type) {
	c.State = state
	c.CharacterAnimationComponent = component.NewCharacterAnimationComponent(c.SpritesheetResource.Spritesheet.Drawables(), .1)
	c.CurrentAction = NewCharacterAction(actionindex.GetActionIndex(state))
	anim := &common.Animation{Name: c.CurrentAction.Name, Frames: c.CurrentAction.Frames}
	c.CharacterAnimationComponent.AddAnimations([]*common.Animation{anim})
	c.CharacterAnimationComponent.AddDefaultAnimation(anim)
	c.CharacterAnimationComponent.CurrentAnimation = anim
}

func (c *CharacterEntity) UUID() string {
	return c.id
}
