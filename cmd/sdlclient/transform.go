package main

import (
	"github.com/go-gl/mathgl/mgl32"
)

var Origin = mgl32.Vec3{0, 0, 0}

type Transform struct {
	position  mgl32.Vec3
	scale     mgl32.Vec3
	direction mgl32.Vec3
	rotation  mgl32.Vec3
}

func NewTransform(position mgl32.Vec3) *Transform {
	return &Transform{
		position:  position,
		scale:     mgl32.Vec3{1, 1, 1},
		direction: Forward,
		rotation:  mgl32.Vec3{0, 0, 0},
	}
}

func (t *Transform) Model() mgl32.Mat4 {
	positionMatrix := mgl32.Translate3D(t.position.X(), t.position.Y(), t.position.Z())
	scaleMatrix := mgl32.Scale3D(t.scale.X(), t.scale.Y(), t.scale.Z())

	rotationX := mgl32.HomogRotate3DX(t.rotation.X())
	rotationY := mgl32.HomogRotate3DY(t.rotation.Y())
	rotationZ := mgl32.HomogRotate3DZ(t.rotation.Z())

	rotationMatrix := rotationZ.Mul4(rotationY).Mul4(rotationX)

	return positionMatrix.Mul4(rotationMatrix).Mul4(scaleMatrix)
}

func (t *Transform) Position() mgl32.Vec3 {
	return t.position
}

func (t *Transform) SetPosition(position mgl32.Vec3) {
	t.position = position
}

func (t *Transform) Scale() mgl32.Vec3 {
	return t.scale
}

func (t *Transform) SetScale(scale mgl32.Vec3) {
	t.scale = scale
}
