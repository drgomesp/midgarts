package character

import (
	"fmt"

	"golang.org/x/text/encoding/charmap"

	"github.com/project-midgard/midgarts/pkg/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/graphics"
)

func LoadCharacterSprite(f *grf.File, jobSpriteID jobspriteid.JobSpriteID) (res *graphics.Sprite, err error) {
	var (
		jobFileName = jobSpriteNameTable[jobSpriteID]
	)

	var decodedFolderA []byte
	if decodedFolderA, err = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC0, 0xCE, 0xB0, 0xA3, 0xC1, 0xB7}); err != nil {
		return nil, err
	}

	var decodedFolderB []byte
	if decodedFolderB, err = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB8, 0xF6, 0xC5, 0xEB}); err != nil {
		return nil, err
	}

	maleFilePath := fmt.Sprintf(`data/sprite/%s/%s/³²/%s_³²`, decodedFolderA, decodedFolderB, jobFileName)
	maleSprite, err := graphics.LoadSprite(f, maleFilePath)
	if err != nil {
		return nil, err
	}

	return maleSprite, nil
}
