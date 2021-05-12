package component

import (
	"github.com/project-midgard/midgarts/pkg/character/actionindex"
	"time"

	"github.com/project-midgard/midgarts/pkg/character/directiontype"
)

type CharacterSpriteRenderInfoComponentFace interface {
	GetCharacterSpriteRenderInfoComponent() *CharacterSpriteRenderInfoComponent
}

type CharacterSpriteRenderInfoComponent struct {
	ActionIndex        actionindex.Type
	AnimationDelay     time.Duration
	AnimationEndsAt    time.Time
	AnimationStartedAt time.Time
	Direction          directiontype.Type
	ForcedDuration     time.Duration
	FPSMultiplier      float64
	IsStandingBy       bool
}

func NewCharacterSpriteRenderInfoComponent() *CharacterSpriteRenderInfoComponent {
	now := time.Now()

	return &CharacterSpriteRenderInfoComponent{
		ActionIndex:        actionindex.Idle,
		AnimationStartedAt: now,
		AnimationEndsAt:    now.Add(time.Millisecond * 100),
		Direction:          directiontype.South,
		ForcedDuration:     0,
		FPSMultiplier:      1.0,
	}
}
