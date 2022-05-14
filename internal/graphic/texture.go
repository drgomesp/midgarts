package graphic

import (
	"fmt"
	"github.com/go-gl/gl/v3.3-core/gl"
	"image"
	"image/draw"
	_ "image/png"
)

type Texture struct {
	handle                        uint32
	path                          string
	width, height, internalFormat int32
	format, formatType            uint32
	magFilter, minFilter          int32
	wrapS, wrapT                  int32
}

func (t *Texture) Bind(unit uint32) {
	gl.ActiveTexture(gl.TEXTURE0 + unit)
	gl.BindTexture(gl.TEXTURE_2D, t.handle)
}

type TextureProvider interface {
	NewTextureFromRGBA(rgba *UniqueRGBA) (tex *Texture, err error)
}

func NewTextureFromRGBA(rgba *UniqueRGBA) (tex *Texture, err error) {
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return nil, fmt.Errorf("unsupported stride")
	}

	draw.Draw(rgba, rgba.Bounds(), rgba, image.Point{}, draw.Src)

	tex = &Texture{
		width:          int32(rgba.Rect.Size().X),
		height:         int32(rgba.Rect.Size().Y),
		internalFormat: gl.RGBA8,
		format:         gl.RGBA,
		formatType:     gl.UNSIGNED_BYTE,
		magFilter:      gl.NEAREST,
		minFilter:      gl.NEAREST,
		wrapS:          gl.CLAMP_TO_BORDER,
		wrapT:          gl.CLAMP_TO_BORDER,
	}

	gl.GenTextures(1, &tex.handle)
	gl.BindTexture(gl.TEXTURE_2D, tex.handle)

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
