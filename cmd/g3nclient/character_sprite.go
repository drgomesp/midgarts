package main

import (
	"fmt"

	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"

	"github.com/g3n/engine/graphic"

	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"

	"github.com/g3n/engine/texture"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
)

type CharacterSprite struct {
	Path string

	actFile  *act.ActionFile
	sprFile  *spr.SpriteFile
	textures []*texture.Texture2D

	bodyMaterial material.IMaterial
	bodySprite   *graphic.Sprite
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

	var filePath string
	if character.Male == gender {
		filePath = character.MaleFilePathf
	} else {
		filePath = character.FemaleFilePathf
	}

	jobFileName := character.JobSpriteNameTable[jobSpriteID]
	path := fmt.Sprintf(filePath, "ÀÎ°£Á·", "¸öÅë", jobFileName)

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

	defaultTexture := textures[0]
	mat := material.NewStandard(math32.NewColor("white"))
	mat.AddTexture(defaultTexture)

	return &CharacterSprite{
		Path: path,

		actFile:  actFile,
		sprFile:  sprFile,
		textures: textures,

		bodyMaterial: mat,
		bodySprite: graphic.NewSprite(
			float32(defaultTexture.Width()),
			float32(defaultTexture.Height()),
			mat,
		),
	}, nil
}

func (s *CharacterSprite) GetTextureAtIndex(i int32) *texture.Texture2D {
	return s.textures[i]
}
