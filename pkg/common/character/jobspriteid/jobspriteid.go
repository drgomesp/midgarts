package jobspriteid

import (
	"log"

	"github.com/project-midgard/midgarts/pkg/common/character/jobid"
)

type Type int

const (
	Novice    Type = 0
	Swordsman Type = 1
	Magician  Type = 2
	Archer    Type = 3
	Alcolyte  Type = 4
	Merchant  Type = 5
	Thief     Type = 6
	Monk      Type = 15
	MonkH     Type = 4016
)

func GetJobSpriteID(jid jobid.Type) (t Type) {
	switch jid {
	case jobid.Archer:
		return Archer
	case jobid.Merchant:
		return Merchant
	case jobid.Monk:
		return Monk
	case jobid.Thief:
		return Thief
	case jobid.MonkH:
		return MonkH
	default:
		log.Fatalf("jobid '%v' not supported", jid)
	}

	return
}
