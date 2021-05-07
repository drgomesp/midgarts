package jobspriteid

import (
	"log"

	"github.com/project-midgard/midgarts/pkg/character/jobid"
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

func GetJobSpriteID(jid jobid.Type, isMounted bool) (t Type) {
	switch jid {
	case jobid.Archer:
		return Archer
	case jobid.Monk:
		return Monk
	case jobid.Assassin:
		return Assassin
	case jobid.Swordsman:
		return Swordsman
	case jobid.Alchemist:
		return Alchemist
	case jobid.Knight:
		if isMounted {
			return Knight2
		} else {
			return Knight
		}
	case jobid.Crusader:
		if isMounted {
			return Crusader2
		} else {
			return Crusader
		}
	//case jobid.Thief:
	//	return Thief
	//case jobid.MonkH:
	//	return MonkH

	default:
		log.Fatalf("jobid '%v' not supported", jid)
	}

	return
}

func (j Type) String() string {
	switch j {
	case Novice:
		return "Novice"
	case Swordsman:
		return "Swordsman"
	case Magician:
		return "Magician"
	case Archer:
		return "Archer"
	case Alcolyte:
		return "Alcolyte"
	case Merchant:
		return "Merchant"
	case Thief:
		return "Thief"
	case Knight:
		return "Knight"
	case Priest:
		return "Priest"
	case Wizard:
		return "Wizard"
	case Blacksmith:
		return "Blacksmith"
	case Hunter:
		return "Hunter"
	case Assassin:
		return "Assassin"
	case Knight2:
		return "Knight2"
	case Crusader:
		return "Crusader"
	case Monk:
		return "Monk"
	case Sage:
		return "Sage"
	case Rogue:
		return "Rogue"
	case Alchemist:
		return "Alchemist"
	case Bard:
		return "Bard"
	case Dancer:
		return "Dancer"
	case Crusader2:
		return "Crusader2"
	case MonkH:
		return "MonkH"
	default:
		log.Fatalf("unsupported jobspriteid %d\n", j)
	}

	return ""
}

func All() []Type {
	return []Type{
		Novice,
		Swordsman,
		Magician,
		Archer,
		Alcolyte,
		Merchant,
		Thief,
		Knight,
		Priest,
		Wizard,
		Blacksmith,
		Hunter,
		Assassin,
		Knight2,
		Crusader,
		Monk,
		Sage,
		Rogue,
		Alchemist,
		Crusader2,
		MonkH,
	}
}
