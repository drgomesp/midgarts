package actionindex

import (
	"fmt"

	"github.com/project-midgard/midgarts/pkg/common/character/statetype"
)

type Type int

const (
	Idle            Type = 0
	Walking         Type = 8
	Sitting         Type = 16
	PickingItem     Type = 24
	StandBy         Type = 32
	Attacking1      Type = 40
	ReceivingDamage Type = 48
	Freeze1         Type = 56
	Dead            Type = 65
	Freeze2         Type = 72
	Attacking2      Type = 80
	Attacking3      Type = 88
	CastingSpell    Type = 96
)

func GetActionIndex(s statetype.Type) (t Type) {
	switch s {
	case statetype.Walking:
		return Walking
	case statetype.Idle:
		return Idle
	case statetype.Attacking:
		return Attacking1
	default:
		panic(fmt.Sprintf("state type '%v' not supported\n", s))
	}

	return
}

func GetStateType(s Type) (t statetype.Type) {
	switch s {
	case Idle:
		return statetype.Idle
	case Walking:
		return statetype.Walking
	case Attacking1:
	case Attacking2:
	case Attacking3:
		return statetype.Attacking
	default:
		panic(fmt.Sprintf("state type '%v' not supported\n", s))
	}

	return
}
