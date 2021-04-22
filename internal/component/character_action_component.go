package component

import (
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
	"github.com/project-midgard/midgarts/pkg/common/character/actionplaymode"
	"github.com/project-midgard/midgarts/pkg/common/character/directiontype"
	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
)

type CharacterActionComponentFace interface {
	GetCharacterActionComponent() *CharacterActionComponent
}

// CharacterActionComponent defines a component that holds information about character actions,
// such as the action itself (Idle, Walking...), the direction (South, North...) and state.
type CharacterActionComponent struct {
	PlayMode  actionplaymode.Type
	Action    actionindex.Type
	Direction directiontype.Type
	State     statetype.Type
}
