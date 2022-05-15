package renderer

import (
	_ "embed"

	"github.com/project-midgard/midgarts/pkg/graphics"
)

const NumVertexAttributes = 3

var (
	//go:embed shader/sprite.vert
	vertexShader string

	//go:embed shader/sprite.frag
	fragmentShader string
)

type Renderer struct {
	pipeline *graphics.Pipeline
}

func New() *Renderer {
	pipeline := graphics.StartPipeline()

	r := &Renderer{}
	r.pipeline = pipeline

	return r
}
