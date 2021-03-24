package client

import (
	"time"

	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"
)

type CharacterAction struct {
	Name               string
	ActionIndex        actionindex.Type
	Frames             []int
	AnimationStartedAt time.Time
}

func NewCharacterAction(actionIndex actionindex.Type) *CharacterAction {
	characterAction := &CharacterAction{
		Name:        string(actionindex.GetStateType(actionIndex)),
		ActionIndex: actionIndex,
	}

	return characterAction
}

func (a *CharacterAction) SetFrames(frames []int) {
	a.Frames = frames
}
