package main

import (
	"encoding/xml"
	"io"

	"github.com/g3n/engine/graphic"
	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"
	"github.com/g3n/engine/texture"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
)

var BackgroundSpriteMaterial = material.NewStandard(math32.NewColor("white"))

type SubTexture struct {
	Texture *texture.Texture2D
	ID      uint32 `xml:"id,attr"`
	X       int32  `xml:"x,attr"`
	Y       int32  `xml:"y,attr"`
	Width   int32  `xml:"width,attr"`
	Height  int32  `xml:"height,attr"`
}

type Spritesheet struct {
	XMLName     xml.Name `xml:"Spritesheet"`
	Texture     *texture.Texture2D
	ImagePath   string        `xml:"imagePath,attr"`
	SubTextures []*SubTexture `xml:"SubTexture"`
}

func (s Spritesheet) URL() string {
	return s.ImagePath
}

func (s *Spritesheet) SubTexture(i uint32) *SubTexture {
	return s.SubTextures[i]
}

func (s *Spritesheet) SpriteAt(i uint32) *graphic.Sprite {
	sub := s.SubTextures[i]
	BackgroundSpriteMaterial.AddTexture(s.Texture)

	return graphic.NewSprite(
		float32(sub.Width),
		float32(sub.Height),
		BackgroundSpriteMaterial,
	)
}

func LoadSpritesheet(sprFile *spr.SpriteFile, r io.Reader, filePath string) (*Spritesheet, error) {
	spritesheet := new(Spritesheet)
	spritesheet.ImagePath = filePath

	var tmp Spritesheet
	err := xml.NewDecoder(r).Decode(&tmp)
	if err != nil {
		return nil, err
	}

	baseTexture, err := texture.NewTexture2DFromImage(filePath)
	if err != nil {
		return nil, err
	}

	spritesheet.Texture = baseTexture
	spritesheet.SubTextures = make([]*SubTexture, len(sprFile.Frames))

	for i := range sprFile.Frames {
		var sub *SubTexture
		for _, s := range tmp.SubTextures {
			if s.ID == uint32(i) {
				sub = s
			}
		}

		if sub != nil {
			w, h := sub.Width, sub.Height

			spritesheet.SubTextures[sub.ID] = &SubTexture{
				ID:     sub.ID,
				X:      sub.X,
				Y:      sub.Y,
				Width:  w,
				Height: h,
			}
		}
	}

	return spritesheet, nil
}
