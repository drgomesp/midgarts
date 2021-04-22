package directiontype

var DirectionTable = [8]int{6, 5, 4, 3, 2, 1, 0, 7}

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
