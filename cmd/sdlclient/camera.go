package main

import (
	"log"

	"github.com/go-gl/mathgl/mgl32"
)

type Projection int

const (
	Perspective = Projection(iota)
	Orthographic
)

type direction int

var (
	Up      = mgl32.Vec3{0, 1, 0}
	Forward = mgl32.Vec3{0, 0, 1}
)

type Camera struct {
	projection             Projection
	fov, aspect, near, far float32

	position mgl32.Vec3
}

func NewPerspectiveCamera(position mgl32.Vec3, fov, aspect, near, far float32) *Camera {
	return &Camera{
		projection: Perspective,
		aspect:     aspect,
		fov:        fov,
		near:       near,
		far:        far,
		position:   position,
	}
}

func (c *Camera) ViewProjection() (vp mgl32.Mat4) {
	switch c.projection {
	case Perspective:
		return mgl32.
			Perspective(c.fov, c.aspect, c.near, c.far).
			Mul4(mgl32.LookAt(
				c.position.X(),
				c.position.Y(),
				c.position.Z(),
				c.position.X()+Forward.X(),
				c.position.Y()+Forward.Y(),
				c.position.Z()+Forward.Z(),
				Up.X(),
				Up.Y(),
				Up.Z(),
			))
	default:
		log.Fatalf("'%v' projection not supported", c.projection)
	}

	return
}
