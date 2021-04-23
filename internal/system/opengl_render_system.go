package system

import (
	"github.com/EngoEngine/ecs"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/project-midgard/midgarts/internal/opengl"
	"github.com/project-midgard/midgarts/internal/system/rendercmd"
	"github.com/project-midgard/midgarts/pkg/camera"
	"github.com/project-midgard/midgarts/pkg/graphic"
)

type RenderCommands struct {
	sprite []rendercmd.SpriteRenderCommand
}

// OpenGLRenderSystem defines an OpenGL-based rendering system.
type OpenGLRenderSystem struct {
	gls            *opengl.State
	cam            *camera.Camera
	renderCommands *RenderCommands
}

func NewOpenGLRenderSystem(gls *opengl.State, cam *camera.Camera, commands *RenderCommands) *OpenGLRenderSystem {
	return &OpenGLRenderSystem{
		gls:            gls,
		cam:            cam,
		renderCommands: commands,
	}
}

func (s *OpenGLRenderSystem) Update(dt float32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(s.gls.Program().ID())

	// 2D Sprites
	{
		for _, cmd := range s.renderCommands.sprite {
			if cmd.FlipVertically {
				cmd.Size = mgl32.Vec2{-cmd.Size.X(), cmd.Size.Y()}
			}

			sprite := graphic.NewSprite(cmd.Size.X(), cmd.Size.Y(), cmd.Texture)
			sprite.SetPosition(mgl32.Vec3{cmd.Position.X(), cmd.Position.Y(), cmd.Position.Z()})
			sprite.Texture.Bind(0)

			view := s.cam.ViewMatrix()
			viewu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("view\x00"))
			gl.UniformMatrix4fv(viewu, 1, false, &view[0])

			model := sprite.Model()
			modelu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("model\x00"))
			gl.UniformMatrix4fv(modelu, 1, false, &model[0])

			projection := s.cam.ProjectionMatrix()
			projectionu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("projection\x00"))
			gl.UniformMatrix4fv(projectionu, 1, false, &projection[0])

			sizeu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("size\x00"))
			gl.Uniform2fv(sizeu, 1, &cmd.Size[0])

			offsetu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("offset\x00"))
			gl.Uniform2fv(offsetu, 1, &cmd.Offset[0])

			iden := mgl32.Ident4()
			rotation := iden.Mul4(mgl32.HomogRotate3D(0, graphic.Forward))
			rotationu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("rotation\x00"))
			gl.UniformMatrix4fv(rotationu, 1, false, &rotation[0])

			sprite.Render(s.gls)
		}
	}
}

func (s *OpenGLRenderSystem) Remove(e ecs.BasicEntity) {
	panic("implement me")
}
