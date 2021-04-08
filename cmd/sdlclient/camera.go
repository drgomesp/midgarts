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
	transform              *Transform
	projection             Projection
	fov, aspect, near, far float32
	projectionMatrix       mgl32.Mat4
}

func NewPerspectiveCamera(position mgl32.Vec3, fov, aspect, near, far float32) *Camera {
	return &Camera{
		transform:  NewTransform(position),
		projection: Perspective,
		aspect:     aspect,
		fov:        fov,
		near:       near,
		far:        far,
	}
}

func NewOrthographicCamera(position mgl32.Vec3, left, right, bottom, top float32) *Camera {
	return &Camera{
		projection: Orthographic,
	}
}

func (c *Camera) ViewProjectionMatrix() (vp mgl32.Mat4) {
	switch c.projection {
	case Perspective:
		return mgl32.
			Perspective(c.fov, c.aspect, c.near, c.far).
			Mul4(mgl32.LookAt(
				c.transform.Position().X(),
				c.transform.Position().Y(),
				c.transform.Position().Z(),
				c.transform.Position().X()+Forward.X(),
				c.transform.Position().Y()+Forward.Y(),
				c.transform.Position().Z()+Forward.Z(),
				Up.X(),
				Up.Y(),
				Up.Z(),
			))
	case Orthographic:
		return mgl32.
			Ortho(0, 1000, 0, 1000, 01, 100).
			Mul4(mgl32.LookAt(
				c.transform.Position().X(),
				c.transform.Position().Y(),
				c.transform.Position().Z(),
				c.transform.Position().X()+Forward.X(),
				c.transform.Position().Y()+Forward.Y(),
				c.transform.Position().Z()+Forward.Z(),
				Up.X(),
				Up.Y(),
				Up.Z(),
			))
	default:
		log.Fatalf("'%v' projection not supported", c.projection)
	}

	return
}

func (c *Camera) SetPosition(position mgl32.Vec3) {
	c.transform.SetPosition(position)
}
