package graphics

import (
	"fmt"

	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
	"golang.org/x/text/encoding/charmap"
)

type Sprite struct {
	act *act.ActionFile
	spr *spr.SpriteFile

	path     string
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
		act:      actFile,
		spr:      sprFile,
		path:     path,
		textures: make([]*common.Texture, 0),
	}

	for i := range sprite.spr.Frames {
		img := sprFile.ImageAt(i)
		if img != nil {
			tex := common.NewTextureSingle(common.NewImageObject(common.ImageToNRGBA(img, img.Bounds().Max.X, img.Bounds().Max.Y)))
			sprite.textures = append(sprite.textures, &tex)
		}
	}

	return
}

func LoadCharacterSprite(f *grf.File, gender character.GenderType, jobSpriteID jobspriteid.Type) (sprite *CharacterSprite, err error) {
	var (
		jobFileName = character.JobSpriteNameTable[jobSpriteID]
	)

	if "" == jobFileName {
		return nil, fmt.Errorf("unsupported jobSpriteID %v", jobSpriteID)
	}

	var decodedFolderA []byte
	if decodedFolderA, err = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC0, 0xCE, 0xB0, 0xA3, 0xC1, 0xB7}); err != nil {
		return nil, err
	}

	var decodedFolderB []byte
	if decodedFolderB, err = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB8, 0xF6, 0xC5, 0xEB}); err != nil {
		return nil, err
	}

	var filePath string
	if character.Male == gender {
		filePath = maleFilePathf
	} else {
		filePath = femaleFilePathf
	}

	bodySprite, err := LoadSprite(f, fmt.Sprintf(filePath, decodedFolderA, decodedFolderB, jobFileName))
	if err != nil {
		return nil, err
	}

	return NewCharacterSprite(gender, bodySprite), nil
}

func LoadMonsterSprite(bodySprite *Sprite) (*CharacterSprite, error) {
	return NewMonsterSprite(bodySprite), nil
}

func (s *Sprite) GetTextureAtIndex(i int32) *common.Texture {
	return s.textures[i]
}
