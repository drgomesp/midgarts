package camera

type Projection int

const (
	Perspective Projection = iota
	Orthographic
)

const (
	Yaw   = float32(270.0)
	Pitch = float32(-60.0)
)
