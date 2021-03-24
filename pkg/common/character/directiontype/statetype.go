package directiontype

type Type uint8

const (
	South Type = iota
	SouthWest
	West
	NorthWest
	North
	NorthEast
	East
	SouthEast
)
