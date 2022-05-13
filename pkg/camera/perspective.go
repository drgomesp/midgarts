package camera

import (
	"math"

	"github.com/drgomesp/midgarts/pkg/graphic"
	"github.com/go-gl/mathgl/mgl32"
)

type PerspectiveCamera struct {
	*graphic.Transform
	target                       *graphic.Transform
	projectionType               Projection
	fov, aspect, near, far       float32
	projectionMatrix, viewMatrix mgl32.Mat4
	front, right, up             mgl32.Vec3
	yaw, pitch                   float32
	distance                     float32
	targetRotation               mgl32.Vec3
	altitude                     float32
}

func NewPerspectiveCamera(fov, aspect, near, far float32) *PerspectiveCamera {
	cam := &PerspectiveCamera{
		Transform:      graphic.NewTransform(mgl32.Vec3{0, 40.0, 0}),
		projectionType: Perspective,
		aspect:         aspect,
		fov:            fov,
		near:           near,
		far:            far,
		distance:       30.0,
		altitude:       50.0,
		front:          mgl32.Vec3{0, 0, -1},
	}

	cam.projectionMatrix = mgl32.Perspective(cam.fov, cam.aspect, cam.near, cam.far)

	return cam
}

func (c *PerspectiveCamera) createViewMatrix() mgl32.Mat4 {
	return mgl32.LookAt(
		c.Position().X(),
		c.Position().Y(),
		c.Position().Z(),
		c.Position().X()+graphic.Forward.X(),
		c.Position().Y()+graphic.Forward.Y(),
		c.Position().Z()+graphic.Forward.Z(),
		graphic.Up.X(),
		graphic.Up.Y(),
		graphic.Up.Z(),
	)
}

func (c *PerspectiveCamera) ViewMatrix() mgl32.Mat4 {
	c.viewMatrix = c.createViewMatrix()
	return c.viewMatrix
}

func (c PerspectiveCamera) ProjectionMatrix() mgl32.Mat4 {
	return c.projectionMatrix
}

func (c *PerspectiveCamera) ResetAngleAndY(windowWidth, windowHeight uint32) {
	c.yaw = Yaw
	c.pitch = Pitch
	c.SetY(40)
	c.Rotate(c.pitch, c.yaw)
	c.UpdateVisibleZRange(windowWidth, windowHeight)
}

func (c *PerspectiveCamera) SetY(y float32) {
	c.SetPosition(mgl32.Vec3{c.Position().X(), y, c.Position().Z() - 32})
}

func (c *PerspectiveCamera) Rotate(yaw float32, pitch float32) {
	c.front = mgl32.Vec3{
		float32(math.Cos(float64(mgl32.DegToRad(pitch))) * math.Cos(float64(mgl32.DegToRad(yaw)))),
		float32(math.Cos(float64(mgl32.DegToRad(pitch)))),
		float32(math.Cos(float64(mgl32.DegToRad(pitch))) * math.Sin(float64(mgl32.DegToRad(yaw)))),
	}

	c.right = c.front.Cross(graphic.Up)
	c.up = c.right.Cross(c.front)
}

func (c *PerspectiveCamera) UpdateVisibleZRange(width uint32, height uint32) {
	view := c.createViewMatrix()
	_ = view
}
