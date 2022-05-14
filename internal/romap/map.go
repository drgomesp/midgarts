package romap

type CellType byte

const (
	CellTypeNone     = CellType(1 << 0)
	CellTypeWalkable = CellType(1 << 1)
	CellTypeWater    = CellType(1 << 2)
	CellTypeSnipable = CellType(1 << 3)
)
