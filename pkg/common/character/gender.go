package character

type GenderType byte

const (
	Male GenderType = iota
	Female
)

func (t GenderType) String() string {
	if t == Male {
		return "m"
	} else {
		return "f"
	}
}
