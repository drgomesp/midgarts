package graphic

import (
	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"
)

type Graphic struct {
	*Transform
	*Geometry
	mode uint32
}

func NewGraphic(mode uint32, geom *Geometry) *Graphic {
	return &Graphic{
		Transform: NewTransform(Origin),
		Geometry:  geom,
		mode:      mode,
	}
}

func (g *Graphic) Render(gls *opengl.State, cam *Camera) {
}
