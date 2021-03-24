package actionplaymode

type Type int

const (
	Repeat Type = iota
	PlayThenHold
	Once
	Reverse
	FixFrame
)
