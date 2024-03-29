package component

import (
	"github.com/project-midgard/midgarts/internal/character/actionindex"
	"github.com/project-midgard/midgarts/internal/character/directiontype"
	"time"
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
