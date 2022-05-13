package rendercmd

import (
	"github.com/drgomesp/midgarts/pkg/graphic"
	"github.com/go-gl/mathgl/mgl32"
)

type SpriteRenderCommand struct {
	Scale           [2]float32
	Size            mgl32.Vec2
	Position        mgl32.Vec3
	Offset          mgl32.Vec2
	RotationRadians float32
	Texture         *graphic.Texture
	FlipVertically  bool
}
