package graphics

import (
	"fmt"

	"github.com/EngoEngine/engo/common"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"golang.org/x/text/encoding/charmap"
)

type SpritesheetResource struct {
	Spritesheet *common.Spritesheet
}

func NewSpritesheetResource(spritesheet *common.Spritesheet) *SpritesheetResource {
	return &SpritesheetResource{
		Spritesheet: spritesheet,
	}
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
		filePath = character.MaleFilePathf
	} else {
		filePath = character.FemaleFilePathf
	}

	bodySprite, err := LoadSprite(f, fmt.Sprintf(filePath, decodedFolderA, decodedFolderB, jobFileName))
	if err != nil {
		return nil, err
	}

	return NewCharacterSprite(gender, bodySprite), nil
}
