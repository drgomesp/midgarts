package graphic

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

var (
	Origin  = mgl32.Vec3{0, 0, 0}
	Up      = mgl32.Vec3{0, 1, 0}
	Forward = mgl32.Vec3{0, 0, 1}
)

type Projection int

const (
	Perspective = Projection(iota)
	Orthographic
)

type Camera struct {
	*Transform
	projection               Projection
	fov, aspect, near, far   float32
	left, right, bottom, top float32
	projectionMatrix         mgl32.Mat4
}

func NewPerspectiveCamera(fov, aspect, near, far float32) *Camera {
	return &Camera{
		Transform:  NewTransform(mgl32.Vec3{0, 0, -5}),
		projection: Perspective,
		aspect:     aspect,
		fov:        fov,
		near:       near,
		far:        far,
	}
}

func NewOrthographicCamera(left, right, bottom, top float32) *Camera {
	return &Camera{
		Transform:  NewTransform(mgl32.Vec3{0, 0, 0}),
		left:       left,
		right:      right,
		bottom:     bottom,
		top:        top,
		projection: Orthographic,
	}
}

func (c *Camera) ViewProjectionMatrix() (vp mgl32.Mat4) {
	switch c.projection {
	case Perspective:
		return mgl32.
			Perspective(c.fov, c.aspect, c.near, c.far).
			Mul4(mgl32.LookAt(
				c.Position().X(),
				c.Position().Y(),
				c.Position().Z(),
				c.Position().X()+Forward.X(),
				c.Position().Y()+Forward.Y(),
				c.Position().Z()+Forward.Z(),
				Up.X(),
				Up.Y(),
				Up.Z(),
			))
	case Orthographic:
		return mgl32.
			Ortho2D(c.left, c.right, c.bottom, c.top).
			Mul4(mgl32.LookAt(
				c.Position().X(),
				c.Position().Y(),
				c.Position().Z(),
				c.Position().X()+Forward.X(),
				c.Position().Y()+Forward.Y(),
				c.Position().Z()+Forward.Z(),
				Up.X(),
				Up.Y(),
				Up.Z(),
			))
	default:
		log.Fatalf("'%v' projection not supported", c.projection)
	}

	return
}

func (c *Camera) AspectRatio() float32 {
	return c.aspect
}

func (c *Camera) Latitude() float32 {
	return 1.0
}
