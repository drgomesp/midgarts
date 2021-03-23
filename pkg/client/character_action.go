package client

import (
	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
)

var CharacterActions = map[actionindex.Type]*CharacterAction{
	actionindex.Idle: NewCharacterAction("idle", actionindex.Idle, []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}, nil),
}

type CharacterAction struct {
	Name        string
	ActionIndex actionindex.Type
	Frames      []int

	animationComponent common.AnimationComponent
}

func NewCharacterAction(name string, actionIndex actionindex.Type, frames []int, spritesheet *common.Spritesheet) *CharacterAction {
	characterAction := &CharacterAction{
		Name:        name,
		ActionIndex: actionIndex,
		Frames:      frames,
	}

	return characterAction
}
