package graphics

import (
	"fmt"

	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
)

type Sprite struct {
	Path string

	ActionFile *act.ActionFile
	SpriteFile *spr.SpriteFile

	textures []*common.Texture
}

func LoadSprite(grfFile *grf.File, path string) (sprite *Sprite, err error) {
	var (
		entry   *grf.Entry
		actFile *act.ActionFile
		sprFile *spr.SpriteFile
	)

	if entry, err = grfFile.GetEntry(fmt.Sprintf("%s.act", path)); err != nil {
		return nil, err
	}

	if actFile, err = act.Load(entry.Data); err != nil {
		return nil, err
	}

	if entry, err = grfFile.GetEntry(fmt.Sprintf("%s.spr", path)); err != nil {
		return nil, err
	}

	if sprFile, err = spr.Load(entry.Data); err != nil {
		return nil, err
	}

	sprite = &Sprite{
		ActionFile: actFile,
		SpriteFile: sprFile,
		Path:       path,
		textures:   make([]*common.Texture, 0),
	}

	for i := range sprite.SpriteFile.Frames {
		if img := sprFile.ImageAt(i); img != nil {
			commonImage := common.ImageToNRGBA(img, img.Bounds().Max.X, img.Bounds().Max.Y)
			tex := common.NewTextureSingle(common.NewImageObject(commonImage))
			sprite.textures = append(sprite.textures, &tex)
		}
	}

	return
}

func (s *Sprite) GetTextureAtIndex(i int32) *common.Texture {
	return s.textures[i]
}
