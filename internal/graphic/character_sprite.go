package graphic

import (
	"fmt"
	"log"

	"github.com/project-midgard/midgarts/cmd/sdlclient/opengl"

	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
	"golang.org/x/text/encoding/charmap"
)

type CharacterSprite struct {
	*Sprite

	act    *act.ActionFile
	spr    *spr.SpriteFile
	Gender character.GenderType
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
		filePath = fmt.Sprintf(character.MaleFilePathf, decodedFolderA, decodedFolderB, jobFileName)
	} else {
		filePath = fmt.Sprintf(character.FemaleFilePathf, decodedFolderA, decodedFolderB, jobFileName)
	}

	var entry *grf.Entry
	if entry, err = f.GetEntry(fmt.Sprintf("%s.act", filePath)); err != nil {
		return nil, err
	}

	var actFile *act.ActionFile
	if actFile, err = act.Load(entry.Data); err != nil {
		return nil, err
	}

	if entry, err = f.GetEntry(fmt.Sprintf("%s.spr", filePath)); err != nil {
		return nil, err
	}

	var sprFile *spr.SpriteFile
	if sprFile, err = spr.Load(entry.Data); err != nil {
		return nil, err
	}

	bodySpriteImage := sprFile.ImageAt(0)
	bodySpriteTexture, err := NewTextureFromImage(bodySpriteImage)
	bodySpriteTexture.Bind(0)

	if err != nil {
		log.Fatal(err)
	}

	w := float32(bodySpriteImage.Rect.Size().X)
	h := float32(bodySpriteImage.Rect.Size().Y)

	return &CharacterSprite{
		act:    actFile,
		spr:    sprFile,
		Sprite: NewSprite(w, h, bodySpriteTexture),
		Gender: gender,
	}, nil

}

func (s *CharacterSprite) Render(gls *opengl.State, cam *Camera) {
	currentFrame := 0

	bodySpriteImage := s.spr.ImageAt(currentFrame)
	bodySpriteTexture, err := NewTextureFromImage(bodySpriteImage)
	if err != nil {
		log.Fatal(err)
	}

	w := float32(bodySpriteImage.Rect.Size().X)
	h := float32(bodySpriteImage.Rect.Size().Y)

	s.Sprite = NewSprite(w, h, bodySpriteTexture)
	s.Sprite.Texture.Bind(0)

	s.Sprite.Render(gls, cam)
}
