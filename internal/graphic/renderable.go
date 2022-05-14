package graphic

import "github.com/go-gl/mathgl/mgl32"

type RenderInfo interface {
	ViewMatrix() mgl32.Mat4
	ProjectionMatrix() mgl32.Mat4
}
