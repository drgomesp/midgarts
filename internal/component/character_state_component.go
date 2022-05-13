package component

import (
	"github.com/drgomesp/midgarts/pkg/character/actionplaymode"
	"github.com/drgomesp/midgarts/pkg/character/statetype"
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
