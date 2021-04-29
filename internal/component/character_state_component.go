package component

import (
	actionplaymode2 "github.com/project-midgard/midgarts/pkg/character/actionplaymode"
	statetype2 "github.com/project-midgard/midgarts/pkg/character/statetype"
)

type CharacterStateComponentFace interface {
	GetCharacterStateComponent() *CharacterStateComponent
}

// CharacterStateComponent defines a component that holds information about character state,
// such as the action (Idle, Walking...), the direction (South, North...) and state.
type CharacterStateComponent struct {
	PlayMode             actionplaymode2.Type
	PreviousState, State statetype2.Type
}
