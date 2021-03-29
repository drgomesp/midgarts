package character

import (
	"fmt"
	"log"

	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/act"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"golang.org/x/text/encoding/charmap"
)

var (
	EncodedDirectoryA = []byte{0xC0, 0xCE, 0xB0, 0xA3, 0xC1, 0xB7}
	EncodedDirectoryB = []byte{0xB8, 0xF6, 0xC5, 0xEB}

	MaleFilePathf   = "data/sprite/%s/%s/³²/%s_³²"
	FemaleFilePathf = "data/sprite/%s/%s/¿©/%s_¿©"
)

func LoadCharacterActionFile(f *grf.File, gender GenderType, jobSpriteID jobspriteid.Type) *act.ActionFile {
	var err error
	path := BuildSpriteFilePath(gender, jobSpriteID)
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

func BuildSpriteFilePath(gender GenderType, jobSpriteID jobspriteid.Type) string {
	var err error
	jobFileName := JobSpriteNameTable[jobSpriteID]

	var decodedFolderA []byte
	if decodedFolderA, err = charmap.Windows1252.NewDecoder().Bytes(EncodedDirectoryA); err != nil {
		log.Fatal(err)
	}

	var decodedFolderB []byte
	if decodedFolderB, err = charmap.Windows1252.NewDecoder().Bytes(EncodedDirectoryB); err != nil {
		log.Fatal(err)
	}

	var filePath string
	if Male == gender {
		filePath = MaleFilePathf
	} else {
		filePath = FemaleFilePathf
	}

	return fmt.Sprintf(filePath, decodedFolderA, decodedFolderB, jobFileName)
}
