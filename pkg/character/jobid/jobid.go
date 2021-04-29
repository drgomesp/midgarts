package jobid

type Type int

const (
	Crusader = Type(iota)
	Swordsman
	Archer
	Ranger
	Assassin
	Rogue
	Knight
	Wizard
	Sage
	Alchemist
	Blacksmith
	Priest
	Monk
	Gunslinger
)

func (t Type) String() string {
	switch t {
	case Crusader:
		return "Crusader"
	case Swordsman:
		return "Swordsman"
	case Archer:
		return "Archer"
	case Ranger:
		return "Ranger"
	case Assassin:
		return "Assassin"
	case Rogue:
		return "Rogue"
	case Knight:
		return "Knight"
	case Wizard:
		return "Wizard"
	case Sage:
		return "Sage"
	case Alchemist:
		return "Alchemist"
	case Blacksmith:
		return "Blacksmith"
	case Priest:
		return "Priest"
	case Monk:
		return "Monk"
	case Gunslinger:
		return "Gunslinger"
	default:
		return ""
	}
}
