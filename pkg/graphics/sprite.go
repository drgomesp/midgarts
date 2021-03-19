package graphics

import (
	"fmt"

	"github.com/EngoEngine/engo/common"

	"github.com/project-midgard/midgarts/pkg/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/fileformat/spr"
)

type Sprite struct {
	act *act.ActionFile
	spr *spr.SpriteFile

	path     string
	Textures []*common.Texture
}

func LoadSprite(grfFile *grf.File, path string) (res *Sprite, err error) {
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

	res = &Sprite{
		act:  actFile,
		spr:  sprFile,
		path: path,
	}

	for i, frame := range res.spr.Frames {
		if frame.SpriteType == spr.FileTypePAL {
			img := sprFile.ImageAt(i)
			tex := common.NewTextureSingle(common.NewImageObject(common.ImageToNRGBA(img, img.Bounds().Max.X, img.Bounds().Max.Y)))
			res.Textures = append(res.Textures, &tex)
		}
	}

	return
}
