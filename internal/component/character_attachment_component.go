package component

import (
	"fmt"
	character2 "github.com/project-midgard/midgarts/pkg/character"
	jobspriteid2 "github.com/project-midgard/midgarts/pkg/character/jobspriteid"
	act2 "github.com/project-midgard/midgarts/pkg/fileformat/act"
	grf2 "github.com/project-midgard/midgarts/pkg/fileformat/grf"
	spr2 "github.com/project-midgard/midgarts/pkg/fileformat/spr"
	"log"

	"github.com/pkg/errors"
	"golang.org/x/text/encoding/charmap"
)

type CharacterAttachmentComponentFace interface {
	GetCharacterAttachmentComponent() *CharacterAttachmentComponent
}

// CharacterAttachmentComponent defines a component that holds state about character
// character attachments (shadow, body, head...).
type CharacterAttachmentComponent struct {
	Files [character2.NumAttachments]struct {
		ACT *act2.ActionFile
		SPR *spr2.SpriteFile
	}
}

func NewCharacterAttachmentComponent(
	f *grf2.File,
	gender character2.GenderType,
	jobSpriteID jobspriteid2.Type,
	headIndex int,
) (*CharacterAttachmentComponent, error) {
	jobFileName := character2.JobSpriteNameTable[jobSpriteID]
	if "" == jobFileName {
		return nil, fmt.Errorf("unsupported jobSpriteID: %v", jobSpriteID)
	}

	decodedFolderA, err := getDecodedFolder([]byte{0xC0, 0xCE, 0xB0, 0xA3, 0xC1, 0xB7})
	if err != nil {
		return nil, err
	}

	decodedFolderB, err := getDecodedFolder([]byte{0xB8, 0xF6, 0xC5, 0xEB})
	if err != nil {
		return nil, err
	}

	var (
		bodyFilePath   string
		shadowFilePath = "data/sprite/shadow"
		headFilePathf  = "data/sprite/ÀÎ°£Á·/¸Ó¸®Åë/%s/%d_%s"
	)

	if character2.Male == gender {
		bodyFilePath = fmt.Sprintf(character2.MaleFilePathf, decodedFolderA, decodedFolderB, jobFileName)
		headFilePathf = fmt.Sprintf(headFilePathf, "³²", headIndex, "³²")
	} else {
		bodyFilePath = fmt.Sprintf(character2.FemaleFilePathf, decodedFolderA, decodedFolderB, jobFileName)
		headFilePathf = fmt.Sprintf(headFilePathf, "¿©", headIndex, "¿©")
	}

	shadowActFile, shadowSprFile, err := f.GetActionAndSpriteFiles(shadowFilePath)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not load shadow act and spr files (%v, %s)", gender, jobSpriteID))
	}

	bodyActFile, bodySprFile, err := f.GetActionAndSpriteFiles(bodyFilePath)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not load body act and spr files (%v, %s)", gender, jobSpriteID))
	}

	headActFile, headSprFile, err := f.GetActionAndSpriteFiles(headFilePathf)
	if err != nil {
		log.Fatal(errors.Wrapf(err, "could not load head act and spr files (%v, %s)", gender, jobSpriteID))
	}

	return &CharacterAttachmentComponent{[character2.NumAttachments]struct {
		ACT *act2.ActionFile
		SPR *spr2.SpriteFile
	}{
		character2.AttachmentShadow: {
			ACT: shadowActFile,
			SPR: shadowSprFile,
		},
		character2.AttachmentBody: {
			ACT: bodyActFile,
			SPR: bodySprFile,
		},
		character2.AttachmentHead: {
			ACT: headActFile,
			SPR: headSprFile,
		},
	}}, nil
}

func getDecodedFolder(buf []byte) (folder []byte, err error) {
	if folder, err = charmap.Windows1252.NewDecoder().Bytes(buf); err != nil {
		return nil, err
	}

	return folder, nil
}
