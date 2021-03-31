package main

import (
	"fmt"
	"os"

	"github.com/project-midgard/midgarts/pkg/common/character/actionindex"

	"github.com/g3n/engine/graphic"

	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"

	"github.com/g3n/engine/texture"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
)

type CharacterSprite struct {
	Path       string
	ActionFile *act.ActionFile
	SpriteFile *spr.SpriteFile

	spritesheet *Spritesheet
	bodySprite  *graphic.Sprite
}

func LoadCharacterSprite(
	grfFile *grf.File,
	gender character.GenderType,
	jobSpriteID jobspriteid.Type,
) (*CharacterSprite, error) {
	var (
		err     error
		entry   *grf.Entry
		actFile *act.ActionFile
		sprFile *spr.SpriteFile
	)

	path := character.BuildSpriteFilePath(gender, jobSpriteID)
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

	textures := make([]*texture.Texture2D, len(sprFile.Frames))

	for i := range sprFile.Frames {
		if rgba := sprFile.ImageAt(i); rgba != nil {
			textures[i] = texture.NewTexture2DFromRGBA(rgba)
		}
	}

	f, err := os.Open("assets/build/f/4016-1.xml")
	if err != nil {
		return nil, err
	}

	spritesheet, err := LoadSpritesheet(sprFile, f, "assets/build/f/4016-1.png")
	if err != nil {
		return nil, err
	}

	return &CharacterSprite{
		Path:        path,
		ActionFile:  actFile,
		SpriteFile:  sprFile,
		spritesheet: spritesheet,
		bodySprite:  spritesheet.SpriteAt(0),
	}, nil
}

func (s *CharacterSprite) GetBodySprite() *graphic.Sprite {
	return s.bodySprite
}

func (s *CharacterSprite) SetActiveBodySprite(i actionindex.Type) {
	s.bodySprite = s.spritesheet.SpriteAt(int32(i))
}

func (s *CharacterSprite) GetBodySpriteAt(i int32) *graphic.Sprite {
	return s.spritesheet.SpriteAt(i)
}

func (s *CharacterSprite) GetSubTextureAt(i int32) *SubTexture {
	return s.spritesheet.SubTexture(i)
}
