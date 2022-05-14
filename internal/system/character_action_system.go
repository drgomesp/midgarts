package system

import (
	"github.com/EngoEngine/ecs"
	"github.com/project-midgard/midgarts/internal/character/actionindex"
	"github.com/project-midgard/midgarts/internal/character/statetype"
	"github.com/project-midgard/midgarts/internal/component"
	"github.com/project-midgard/midgarts/internal/entity"
	"github.com/project-midgard/midgarts/internal/fileformat/grf"
	"log"
	"strconv"
	"time"
)

type CharacterActionable interface {
	component.CharacterStateComponentFace
	component.CharacterSpriteRenderInfoComponentFace
}

type CharacterActionSystem struct {
	grfFile *grf.File

	characters map[string]*entity.Character
}

func NewCharacterActionSystem(grfFile *grf.File) *CharacterActionSystem {
	return &CharacterActionSystem{
		grfFile,
		map[string]*entity.Character{},
	}
}

func (s *CharacterActionSystem) Add(char *entity.Character) {
	cmp, e := component.NewCharacterAttachmentComponent(s.grfFile, component.CharacterAttachmentComponentConfig{
		Gender:      char.Gender,
		JobSpriteID: char.JobSpriteID,
		HeadIndex:   char.HeadIndex,
	})
	if e != nil {
		log.Fatal(e)
	}
	char.SetCharacterAttachmentComponent(cmp)
	s.characters[strconv.Itoa(int(char.ID()))] = char
}

func (s CharacterActionSystem) AddByInterface(o ecs.Identifier) {
	char := o.(*entity.Character)
	s.Add(char)
}

func (s CharacterActionSystem) Update(dt float32) {
	for _, c := range s.characters {
		now := time.Now()
		previousAnimationHasEnded := now.After(c.AnimationEndsAt)

		stopPreviousAnimation := previousAnimationHasEnded
		if c.PreviousState != statetype.Walking {
			stopPreviousAnimation = true
		}

		c.ActionIndex = actionindex.GetActionIndex(c.State)

		if (c.State != c.PreviousState && c.State != statetype.Idle) ||
			(c.State == statetype.Idle && stopPreviousAnimation) {
			c.AnimationStartedAt = now

			// TODO: treat special case when attacking
			var forcedDuration time.Duration
			c.ForcedDuration = forcedDuration

			c.FPSMultiplier = 1.0
			if c.State == statetype.Walking {
				c.FPSMultiplier = c.MovementSpeed
			}
		}
		c.AnimationEndsAt = now.Add(c.AnimationDelay)
	}
}

func (s CharacterActionSystem) Remove(e ecs.BasicEntity) {
	delete(s.characters, strconv.Itoa(int(e.ID())))
}
