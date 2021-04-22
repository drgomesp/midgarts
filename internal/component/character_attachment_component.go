package component

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobid"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/spr"
	"golang.org/x/text/encoding/charmap"
)

type CharacterAttachmentComponentFace interface {
	GetCharacterAttachmentComponent() *CharacterAttachmentComponent
}

// CharacterAttachmentComponent defines a component that holds state about character
// character attachments (shadow, body, head...).
type CharacterAttachmentComponent struct {
	Files [character.NumAttachments]struct {
		ACT *act.ActionFile
		SPR *spr.SpriteFile
	}
}

func NewCharacterAttachmentComponent(
	f *grf.File,
	gender character.GenderType,
	job jobid.Type,
	headIndex int,
	isMounted bool,
) (*CharacterAttachmentComponent, error) {
	jobSpriteID := jobspriteid.GetJobSpriteID(job, isMounted)
	jobFileName := character.JobSpriteNameTable[jobSpriteID]
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

	if character.Male == gender {
		bodyFilePath = fmt.Sprintf(character.MaleFilePathf, decodedFolderA, decodedFolderB, jobFileName)
		headFilePathf = fmt.Sprintf(headFilePathf, "³²", headIndex, "³²")
	} else {
		bodyFilePath = fmt.Sprintf(character.FemaleFilePathf, decodedFolderA, decodedFolderB, jobFileName)
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

	return &CharacterAttachmentComponent{[6]struct {
		ACT *act.ActionFile
		SPR *spr.SpriteFile
	}{
		character.AttachmentShadow: {
			ACT: shadowActFile,
			SPR: shadowSprFile,
		},
		character.AttachmentBody: {
			ACT: bodyActFile,
			SPR: bodySprFile,
		},
		character.AttachmentHead: {
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
