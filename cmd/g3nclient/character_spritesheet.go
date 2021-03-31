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

type Spritesheet struct {
	XMLName     xml.Name `xml:"Spritesheet"`
	Texture     *texture.Texture2D
	ImagePath   string        `xml:"imagePath,attr"`
	SubTextures []*SubTexture `xml:"SubTexture"`
	Sprites     []*graphic.Sprite
}

func (s Spritesheet) URL() string {
	return s.ImagePath
}

func (s *Spritesheet) SubTexture(i int32) *SubTexture {
	return s.SubTextures[i]
}

func (s *Spritesheet) SpriteAt(i int32) *graphic.Sprite {
	return s.Sprites[i]
}

type SubTexture struct {
	Texture  *texture.Texture2D
	Name     string `xml:"name,attr"`
	X        int32  `xml:"x,attr"`
	Y        int32  `xml:"y,attr"`
	Width    int32  `xml:"width,attr"`
	Height   int32  `xml:"height,attr"`
	FlippedY bool
}

func LoadSpritesheet(sprFile *spr.SpriteFile, r io.Reader, filePath string) (*Spritesheet, error) {
	atlas := new(Spritesheet)
	atlas.ImagePath = filePath

	err := xml.NewDecoder(r).Decode(&atlas)
	if err != nil {
		return nil, err
	}

	baseTexture, err := texture.NewTexture2DFromImage(filePath)
	if err != nil {
		return nil, err
	}

	atlas.Texture = baseTexture
	atlas.Sprites = make([]*graphic.Sprite, len(atlas.SubTextures))

	for i := 0; i < len(atlas.SubTextures); i++ {
		tex := texture.NewTexture2DFromRGBA(sprFile.ImageAt(i))
		mat := material.NewStandard(math32.NewColor("white"))
		mat.AddTexture(tex)
		atlas.SubTextures[i].Texture = tex
		atlas.Sprites[i] = graphic.NewSprite(float32(tex.Width()), float32(tex.Height()), mat)
	}

	return atlas, nil
}
