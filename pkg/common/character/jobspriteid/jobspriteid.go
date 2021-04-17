package jobspriteid

import (
	"log"

	"github.com/project-midgard/midgarts/pkg/common/character/jobid"
)

type Type int

const (
	Novice     = Type(0)
	Swordsman  = Type(1)
	Magician   = Type(2)
	Archer     = Type(3)
	Alcolyte   = Type(4)
	Merchant   = Type(5)
	Thief      = Type(6)
	Knight     = Type(7)
	Priest     = Type(8)
	Wizard     = Type(9)
	Blacksmith = Type(10)
	Hunter     = Type(11)
	Assassin   = Type(12)
	Knight2    = Type(13)
	Crusader   = Type(14)
	Monk       = Type(15)
	Sage       = Type(16)
	Rogue      = Type(17)
	Alchemist  = Type(18)
	Bard       = Type(19)
	Dancer     = Type(20)
	Crusader2  = Type(21)
	MonkH      = Type(4016)
)

func GetJobSpriteID(jid jobid.Type) (t Type) {
	switch jid {
	case jobid.Archer:
		return Archer
	//case jobid.Merchant:
	//	return Merchant
	case jobid.Monk:
		return Monk
	case jobid.Assassin:
		return Assassin
	//case jobid.Thief:
	//	return Thief
	//case jobid.MonkH:
	//	return MonkH

	default:
		log.Fatalf("jobid '%v' not supported", jid)
	}

	return
}
