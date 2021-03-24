package graphics

import (
	"fmt"
	"log"

	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
	"golang.org/x/text/encoding/charmap"
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
		filePath = MaleFilePathf
	} else {
		filePath = FemaleFilePathf
	}

	bodySprite, err := LoadSprite(f, fmt.Sprintf(filePath, decodedFolderA, decodedFolderB, jobFileName))
	if err != nil {
		return nil, err
	}

	return NewCharacterSprite(gender, bodySprite), nil
}

func LoadCharacterActionFile(f *grf.File, gender character.GenderType, jobSpriteID jobspriteid.Type) *act.ActionFile {
	var (
		err         error
		jobFileName = character.JobSpriteNameTable[jobSpriteID]
	)

	if "" == jobFileName {
		log.Fatalf("unsupported jobSpriteID %v", jobSpriteID)
	}

	var decodedFolderA []byte
	if decodedFolderA, err = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC0, 0xCE, 0xB0, 0xA3, 0xC1, 0xB7}); err != nil {
		log.Fatal(err)
	}

	var decodedFolderB []byte
	if decodedFolderB, err = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB8, 0xF6, 0xC5, 0xEB}); err != nil {
		log.Fatal(err)
	}

	var filePath string
	if character.Male == gender {
		filePath = MaleFilePathf
	} else {
		filePath = FemaleFilePathf
	}

	path := fmt.Sprintf(filePath, decodedFolderA, decodedFolderB, jobFileName)
	var entry *grf.Entry
	if entry, err = f.GetEntry(fmt.Sprintf("%s.act", path)); err != nil {
		log.Fatal(err)
	}

	actFile, err := act.Load(entry.Data)
	if err != nil {
		log.Fatal(err)
	}

	return actFile
}

func (s *Sprite) GetTextureAtIndex(i int32) *common.Texture {
	return s.textures[i]
}
