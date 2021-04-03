package main

import (
	"fmt"
	"os"

	"github.com/g3n/engine/material"
	"github.com/g3n/engine/math32"

	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"

	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
)

type CharacterSprite struct {
	Path        string
	ActionFile  *act.ActionFile
	SpriteFile  *spr.SpriteFile
	Spritesheet *Spritesheet
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

	f, err := os.Open("assets/build/f/15-1.xml")
	if err != nil {
		return nil, err
	}

	spritesheet, err := LoadSpritesheet(sprFile, f, "assets/build/f/15-1.png")
	if err != nil {
		return nil, err
	}

	baseMaterial := material.NewStandard(math32.NewColor("white"))
	baseMaterial.AddTexture(spritesheet.Texture)

	return &CharacterSprite{
		Spritesheet: spritesheet,
		Path:        path,
		ActionFile:  actFile,
		SpriteFile:  sprFile,
	}, nil
}
