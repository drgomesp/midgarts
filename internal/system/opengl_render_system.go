package system

import (
	"github.com/EngoEngine/ecs"
	"github.com/davecgh/go-spew/spew"
	"github.com/drgomesp/midgarts/internal/opengl"
	"github.com/drgomesp/midgarts/internal/system/rendercmd"
	"github.com/drgomesp/midgarts/pkg/graphic"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/rs/zerolog/log"
)

type RenderCommands struct {
	sprite []rendercmd.SpriteRenderCommand
}

// OpenGLRenderSystem defines an OpenGL-based rendering system.
type OpenGLRenderSystem struct {
	gls            *opengl.State
	renderInfo     graphic.RenderInfo
	renderCommands *RenderCommands

	// Buffer of reusable sprites
	spritesBuf []*graphic.Sprite
}

func NewOpenGLRenderSystem(gls *opengl.State, renderInfo graphic.RenderInfo, commands *RenderCommands) *OpenGLRenderSystem {
	return &OpenGLRenderSystem{
		gls:            gls,
		renderInfo:     renderInfo,
		renderCommands: commands,
		spritesBuf:     []*graphic.Sprite{},
	}
}

func (s *OpenGLRenderSystem) Update(dt float32) {
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	gl.UseProgram(s.gls.Program().ID())

	// 2D Sprites
	{
		s.EnsureSpritesBufLen(len(s.renderCommands.sprite))

		for i, cmd := range s.renderCommands.sprite {
			log.Trace().Msgf("OpenGLRenderSystem::Update(%v) cmd=(%v)", dt, spew.Sdump(cmd))

			//if cmd.FlipVertically {
			cmd.Size = mgl32.Vec2{cmd.Size.X(), -cmd.Size.Y()}
			//}

			sprite := s.spritesBuf[i]
			sprite.SetBounds(cmd.Size.X(), cmd.Size.Y())
			sprite.SetTexture(cmd.Texture)
			sprite.SetPosition(mgl32.Vec3{cmd.Position.X(), cmd.Position.Y(), cmd.Position.Z()})
			sprite.Texture.Bind(0)

			//view := s.renderInfo.ViewMatrix()
			//viewu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("view\x00"))
			//gl.UniformMatrix4fv(viewu, 1, false, &view[0])

			model := sprite.Model()
			modelu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("model\x00"))
			gl.UniformMatrix4fv(modelu, 1, false, &model[0])

			projection := s.renderInfo.ProjectionMatrix()
			projectionu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("projection\x00"))
			gl.UniformMatrix4fv(projectionu, 1, false, &projection[0])

			sizeu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("size\x00"))
			gl.Uniform2fv(sizeu, 1, &cmd.Size[0])

			offsetu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("offset\x00"))
			gl.Uniform2fv(offsetu, 1, &cmd.Offset[0])

			iden := mgl32.Ident4()
			rotation := iden.Mul4(mgl32.HomogRotate3D(cmd.RotationRadians, graphic.Backwards))
			rotationu := gl.GetUniformLocation(s.gls.Program().ID(), gl.Str("rotation\x00"))
			gl.UniformMatrix4fv(rotationu, 1, false, &rotation[0])

			sprite.Render(s.gls)
		}
	}
}

func (s *OpenGLRenderSystem) Remove(e ecs.BasicEntity) {

}

func (s *OpenGLRenderSystem) EnsureSpritesBufLen(minLen int) {
	s.spritesBuf = ensureSpritesBufferLength(s.spritesBuf, minLen)
}

func ensureSpritesBufferLength(slice []*graphic.Sprite, minLen int) []*graphic.Sprite {
	oldLen := len(slice)

	if cacheOverflow := minLen - oldLen; cacheOverflow <= 0 {
		// no need to resize
		return slice
	}

	if minLen > cap(slice) {
		newSlice := make([]*graphic.Sprite, oldLen, minLen)
		copy(newSlice, slice)
		slice = newSlice
	}

	slice = slice[0:minLen]

	for i := oldLen; i < minLen; i++ {
		slice[i] = graphic.NewSprite(0, 0, nil)
	}
	return slice
}
