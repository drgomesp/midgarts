package graphic

import (
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/project-midgard/midgarts/cmd/sdlclient/array"
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
)

type Sprite struct {
	*Graphic
}

func NewSprite(width, height float32) *Sprite {
	geom := NewGeometry()
	w := width / 2
	h := height / 2

	vb := []float32{
		-w, -h, 0, 1, 1, 1, 0, 0,
		w, -h, 0, 1, 1, 1, 1, 0,
		w, h, 0, 1, 1, 1, 1, 1,
		-w, h, 0, 1, 1, 1, 0, 1,
	}

	indices := array.NewArrayUint32(0, 6)
	indices.Append(0, 1, 2, 0, 2, 3)

	geom.SetIndices(indices)
	geom.AddVBO(
		opengl.NewVBO(vb).
			AddAttribute(opengl.VertexPosition).
			AddAttribute(opengl.VertexColor).
			AddAttribute(opengl.VertexTexCoord),
	)

	return &Sprite{
		Graphic: NewGraphic(gl.TRIANGLES, geom),
	}
}

func (s *Sprite) PreRender(gls *opengl.State, cam *Camera) {
	mvp := cam.ViewProjectionMatrix().Mul4(s.Model())
	mvpUniform := gl.GetUniformLocation(gls.Program().ID(), gl.Str("mvp\x00"))
	gl.UniformMatrix4fv(mvpUniform, 1, false, &mvp[0])
}

func (s *Sprite) Render(gls *opengl.State, cam *Camera) {
	s.Geometry.PreRender(gls, cam)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, uint32(opengl.NumVertexBuffers))

	indices := s.Geometry.Indices()

	if indices.Size() > 0 {
		gl.DrawElements(s.mode, int32(len(s.Indices())), gl.UNSIGNED_INT, gl.Ptr(nil))
	}

	gl.DisableVertexAttribArray(0)
	gl.BindVertexArray(0)
}
