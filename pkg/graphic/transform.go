package graphic

import (
	"github.com/go-gl/mathgl/mgl32"
)

var (
	Origin  = mgl32.Vec3{0, 0, 0}
	Up      = mgl32.Vec3{0, 1, 0}
	Forward = mgl32.Vec3{0, 0, 1}
)

type Transform struct {
	position  mgl32.Vec3
	scale     mgl32.Vec3
	direction mgl32.Vec3
	rotation  mgl32.Quat
}

func NewTransform(position mgl32.Vec3) *Transform {
	return &Transform{
		position:  position,
		scale:     mgl32.Vec3{1, 1, 1},
		direction: Forward,
		rotation:  mgl32.AnglesToQuat(0, 0, 0, mgl32.XYZ),
	}
}

func (t *Transform) Model() mgl32.Mat4 {
	positionMatrix := mgl32.Translate3D(t.position.X(), t.position.Y(), t.position.Z())
	scaleMatrix := mgl32.Scale3D(t.scale.X(), t.scale.Y(), t.scale.Z())

	rotationX := mgl32.HomogRotate3DX(t.rotation.X())
	rotationY := mgl32.HomogRotate3DY(t.rotation.Y())
	rotationZ := mgl32.HomogRotate3DZ(t.rotation.Z())

	rotationMatrix := rotationZ.Mul4(rotationY).Mul4(rotationX)

	model := positionMatrix.Mul4(rotationMatrix).Mul4(scaleMatrix)

	return model
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

func (t *Transform) Rotation() mgl32.Quat {
	return t.rotation
}

func (t *Transform) SetRotation(rotation mgl32.Quat) {
	t.rotation = rotation
}
