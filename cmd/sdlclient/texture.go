package main

import (
	"fmt"
	"image"
	"image/draw"
	_ "image/png"
	"os"

	"github.com/go-gl/gl/v3.3-core/gl"
)

type Texture struct {
	name                          uint32
	path                          string
	img                           *image.RGBA
	width, height, internalFormat int32
	format, formatType            uint32
	magFilter, minFilter          int32
	wrapS, wrapT                  int32
}

func NewTextureFromImage(path string) (tex *Texture, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), img, image.Point{}, draw.Src)

	tex = &Texture{
		path:           path,
		width:          int32(rgba.Rect.Size().X),
		height:         int32(rgba.Rect.Size().Y),
		internalFormat: gl.RGBA8,
		format:         gl.RGBA,
		formatType:     gl.UNSIGNED_BYTE,
		img:            rgba,
		magFilter:      gl.LINEAR,
		minFilter:      gl.LINEAR,
		wrapS:          gl.REPEAT,
		wrapT:          gl.REPEAT,
	}

	gl.GenTextures(1, &tex.name)
	gl.BindTexture(gl.TEXTURE_2D, tex.name)

	gl.GenerateMipmap(gl.TEXTURE_2D)

	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, tex.magFilter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, tex.minFilter)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, tex.wrapS)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, tex.wrapT)

	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		tex.internalFormat,
		tex.width,
		tex.height,
		0,
		tex.format,
		tex.formatType,
		gl.Ptr(&rgba.Pix[0]),
	)

	return
}

func (t *Texture) Bind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.name)
}
