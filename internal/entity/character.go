package entity

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo"
	"github.com/EngoEngine/engo/common"
	clientcomponent "github.com/project-midgard/midgarts/internal/component"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/jobid"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	graphic2 "github.com/project-midgard/midgarts/pkg/graphic"
	uuid "github.com/satori/go.uuid"
)

type Character struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	clientcomponent.CharacterAnimationComponent

	id string

	// spritesheetResource should come from a shared system later, but for now its ok
	// for every character to hold its own spritesheetResource.
	SpritesheetResource *graphic2.SpritesheetResource
	ActionFile          *act.ActionFile
	CurrentAction       *CharacterAction
	PlayMode            actionplaymode.Type

	Gender         character.GenderType
	Job            jobid.Type
	State          statetype.Type
	Direction      directiontype.Type
	TargetPosition engo.Point
}

func NewCharacterEntity(spritesheetResource *graphic2.SpritesheetResource, actFile *act.ActionFile, gender character.GenderType, job jobid.Type) *Character {
	b := ecs.NewBasic()
	return &Character{
		BasicEntity:         b,
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

func (c *Character) SetAction(state statetype.Type) {
	c.State = state
	c.CharacterAnimationComponent = clientcomponent.NewCharacterAnimationComponent(c.SpritesheetResource.Spritesheet.Drawables(), 0.08)
	c.CurrentAction = NewCharacterAction(actionindex.GetActionIndex(state))
	anim := &common.Animation{Name: c.CurrentAction.Name, Frames: c.CurrentAction.Frames}
	c.CharacterAnimationComponent.AddAnimations([]*common.Animation{anim})
	c.CharacterAnimationComponent.AddDefaultAnimation(anim)
	c.CharacterAnimationComponent.CurrentAnimation = anim
}

func (c *Character) UUID() string {
	return c.id
}
