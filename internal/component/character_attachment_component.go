package component

import (
	"fmt"
	"strconv"

	"github.com/pkg/errors"
	"github.com/project-midgard/midgarts/pkg/character"
	"github.com/project-midgard/midgarts/pkg/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/fileformat/grf"
	"golang.org/x/text/encoding/charmap"
)

type CharacterAttachmentComponentFace interface {
	GetCharacterAttachmentComponent() *CharacterAttachmentComponent
}

// CharacterAttachmentComponent defines a component that holds state about character
// character attachments (shadow, body, head...).
type CharacterAttachmentComponent struct {
	Files [character.NumAttachments]grf.ActionSpriteFilePair
}

type CharacterAttachmentComponentConfig struct {
	Gender           character.GenderType
	JobSpriteID      jobspriteid.Type
	HeadIndex        int
	EnableShield     bool
	ShieldSpriteName string // Loaded from GRF
}

func NewCharacterAttachmentComponent(
	f *grf.File,
	conf CharacterAttachmentComponentConfig,
) (*CharacterAttachmentComponent, error) {
	cmp := &CharacterAttachmentComponent{
		Files: [character.NumAttachments]grf.ActionSpriteFilePair{},
	}

	jobFileName := character.JobSpriteNameTable[conf.JobSpriteID]
	if "" == jobFileName {
		return cmp, fmt.Errorf("unsupported jobSpriteID: %v", conf.JobSpriteID)
	}

	decodedFolderA, err := getDecodedFolder([]byte{0xC0, 0xCE, 0xB0, 0xA3, 0xC1, 0xB7})
	if err != nil {
		return cmp, errors.Wrap(err, "unable to decode folder name")
	}

	decodedFolderB, err := getDecodedFolder([]byte{0xB8, 0xF6, 0xC5, 0xEB})
	if err != nil {
		return cmp, errors.Wrap(err, "unable to decode folder name")
	}

	genderPath := "³²"
	if conf.Gender == character.Female {
		genderPath = "¿©"
	}

	cmp.Files[character.AttachmentShadow], err = f.GetActionAndSpriteFiles("data/sprite/shadow")
	if err != nil {
		return cmp, errors.Wrapf(err, "could not load shadow act and spr files (%v, %s)", conf.Gender, conf.JobSpriteID)
	}

	bodyFilePath := "data/sprite/" + decodedFolderA + "/" + decodedFolderB + "/" + genderPath + "/" + jobFileName + "_" + genderPath
	cmp.Files[character.AttachmentBody], err = f.GetActionAndSpriteFiles(bodyFilePath)
	if err != nil {
		return cmp, errors.Wrapf(err, "could not load body act and spr files (%v, %s)", conf.Gender, conf.JobSpriteID)
	}

	headFilePath := "data/sprite/ÀÎ°£Á·/¸Ó¸®Åë/" + genderPath + "/" + strconv.Itoa(conf.HeadIndex) + "_" + genderPath
	cmp.Files[character.AttachmentHead], err = f.GetActionAndSpriteFiles(headFilePath)
	if err != nil {
		return cmp, errors.Wrapf(err, "could not load head act and spr files (%v, %s)", conf.Gender, conf.JobSpriteID)
	}

	if conf.EnableShield {
		if conf.ShieldSpriteName == "" {
			conf.ShieldSpriteName = "°¡µå"
		}
		shieldFilePath := "data/sprite/¹æÆÐ/" + jobFileName + "/" + jobFileName + "_" + genderPath + "_" + conf.ShieldSpriteName
		cmp.Files[character.AttachmentShield], err = f.GetActionAndSpriteFiles(shieldFilePath)
		if err != nil {
			return cmp, errors.Wrapf(err, "could not load shield act and spr files (%v, %s, %s)", conf.Gender, conf.JobSpriteID, conf.ShieldSpriteName)
		}
	}

	return cmp, nil
}

func getDecodedFolder(buf []byte) (string, error) {
	folderNameBytes, err := charmap.Windows1252.NewDecoder().Bytes(buf)
	return string(folderNameBytes), err
}
