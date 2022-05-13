package camera

import (
	"github.com/drgomesp/midgarts/pkg/graphic"
	"github.com/go-gl/mathgl/mgl32"
)

type OrthographicCamera struct {
	*graphic.Transform
	projectionType               Projection
	projectionMatrix, viewMatrix mgl32.Mat4
}

func NewOrthographicCamera(left, right, bottom, top float32) *OrthographicCamera {
	cam := &OrthographicCamera{
		Transform:      graphic.NewTransform(mgl32.Vec3{0, 0, 0}),
		projectionType: Orthographic,
	}

	cam.projectionMatrix = mgl32.Ortho2D(left, right, bottom, top)

	return cam
}

func (c *OrthographicCamera) ViewMatrix() mgl32.Mat4 {
	panic("TODO")
}

func (c OrthographicCamera) ProjectionMatrix() mgl32.Mat4 {
	return c.projectionMatrix
}
