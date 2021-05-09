package actionindex

import (
	"fmt"
	"log"

	"github.com/project-midgard/midgarts/pkg/character/statetype"
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
	case statetype.Attacking:
		return Attacking3
	case statetype.Walking:
		return Walking
	case statetype.Idle:
		return Idle
	case statetype.StandBy:
		return StandBy
	default:
		log.Fatal(fmt.Sprintf("state type '%v' not supported\n", s))
	}

	return
}

func GetStateType(s Type) (t statetype.Type) {
	switch s {
	case Idle:
		return statetype.Idle
	case Walking:
		return statetype.Walking
	case StandBy:
		return statetype.StandBy
	default:
		panic(fmt.Sprintf("state type '%v' not supported\n", s))
	}

	return
}
