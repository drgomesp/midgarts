package client

import (
	"github.com/EngoEngine/ecs"
	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobid"
	uuid "github.com/satori/go.uuid"
)

type Character struct {
	ecs.BasicEntity
	common.RenderComponent
	common.SpaceComponent
	common.AnimationComponent

	id string

	// spritesheet should come from a shared system later, but for now its ok
	// for every character to hold its own spritesheet.
	spritesheet *common.Spritesheet

	gender character.GenderType
	job    jobid.Type

	currentAction *CharacterAction
}

func NewCharacter(spritesheet *common.Spritesheet, gender character.GenderType, job jobid.Type) *Character {
	return &Character{
		BasicEntity: ecs.NewBasic(),
		id:          uuid.NewV4().String(),
		spritesheet: spritesheet,
		gender:      gender,
		job:         job,
	}
}

func (c *Character) SetCurrentAction(act *CharacterAction) {
	c.currentAction = act

	c.AnimationComponent = common.NewAnimationComponent(c.spritesheet.Drawables(), .09)

	animationAction0 := &common.Animation{Name: act.Name, Frames: act.Frames}

	c.AnimationComponent.AddAnimations([]*common.Animation{animationAction0})
	c.AnimationComponent.AddDefaultAnimation(animationAction0)
}

func (c *Character) Spritesheet() *common.Spritesheet {
	return c.spritesheet
}
