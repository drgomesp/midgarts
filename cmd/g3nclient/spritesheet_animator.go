package main

import (
	"log"
	"time"

	"github.com/g3n/engine/texture"

	"github.com/g3n/engine/gls"
)

type Animator struct {
	Spritesheet                        *Spritesheet
	Texture                            *texture.Texture2D
	CurrentFrame                       uint32
	OffsetX, OffsetY, RepeatX, RepeatY float32
}

func NewAnimator(spritesheet *Spritesheet) *Animator {
	a := &Animator{
		Spritesheet:  spritesheet,
		Texture:      spritesheet.Texture,
		CurrentFrame: 0,
	}

	a.Texture.SetWrapS(gls.REPEAT)
	a.Texture.SetWrapT(gls.REPEAT)
	a.Texture.SetMagFilter(gls.LINEAR)

	a.Update(time.Now())

	return a
}

func (a *Animator) Update(now time.Time) {
	sub := a.Spritesheet.SubTexture(a.CurrentFrame)
	a.OffsetX = float32(sub.X) / float32(a.Texture.Width())
	a.OffsetY = float32(sub.Y) / float32(a.Texture.Height())
	a.RepeatX = 1 / float32(a.Texture.Width()/int(sub.Width))
	a.RepeatY = 1 / float32(a.Texture.Height()/int(sub.Height))

	a.Texture.SetOffset(a.OffsetX, a.OffsetY)
	a.Texture.SetRepeat(a.RepeatX, a.RepeatY)

	log.Printf("repeat=(%v,%v)\n", a.RepeatX, a.RepeatY)
}
