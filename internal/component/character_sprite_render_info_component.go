package component

import (
	character2 "github.com/project-midgard/midgarts/pkg/character"
	actionindex2 "github.com/project-midgard/midgarts/pkg/character/actionindex"
	directiontype2 "github.com/project-midgard/midgarts/pkg/character/directiontype"
	"time"
)

type CharacterSpriteRenderInfoComponentFace interface {
	GetCharacterSpriteRenderInfoComponent() *CharacterSpriteRenderInfoComponent
}

type CharacterSpriteRenderInfoComponent struct {
	ActionIndex        actionindex2.Type
	AnimationEndsAt    time.Time
	AnimationStartedAt time.Time
	Direction          directiontype2.Type
	ForcedDuration     time.Duration
	FPSMultiplier      float64
	AttachmentType     character2.AttachmentType
}

func NewCharacterSpriteRenderInfoComponent() *CharacterSpriteRenderInfoComponent {
	now := time.Now()

	return &CharacterSpriteRenderInfoComponent{
		ActionIndex:        actionindex2.Idle,
		AnimationStartedAt: now,
		AnimationEndsAt:    now.Add(time.Millisecond * 100),
		Direction:          directiontype2.South,
		ForcedDuration:     0,
		FPSMultiplier:      1.0,
	}
}
