package component

import (
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
)

type CharacterStateComponentFace interface {
	GetCharacterStateComponent() *CharacterStateComponent
}

// CharacterStateComponent defines a component that holds information about character state,
// such as the action (Idle, Walking...), the direction (South, North...) and state.
type CharacterStateComponent struct {
	PlayMode             actionplaymode.Type
	PreviousState, State statetype.Type
}
