package character

import (
	"fmt"
	jobspriteid2 "github.com/project-midgard/midgarts/pkg/character/jobspriteid"
	act2 "github.com/project-midgard/midgarts/pkg/fileformat/act"
	grf2 "github.com/project-midgard/midgarts/pkg/fileformat/grf"
	"log"

	"golang.org/x/text/encoding/charmap"
)

var (
	EncodedDirectoryA = []byte{0xC0, 0xCE, 0xB0, 0xA3, 0xC1, 0xB7}
	EncodedDirectoryB = []byte{0xB8, 0xF6, 0xC5, 0xEB}

	MaleFilePathf   = "data/sprite/%s/%s/³²/%s_³²"
	FemaleFilePathf = "data/sprite/%s/%s/¿©/%s_¿©"
)

func LoadCharacterActionFile(f *grf2.File, gender GenderType, jobSpriteID jobspriteid2.Type) *act2.ActionFile {
	var err error
	path := BuildSpriteFilePath(gender, jobSpriteID)
	var entry *grf2.Entry
	if entry, err = f.GetEntry(fmt.Sprintf("%s.act", path)); err != nil {
		log.Fatal(err)
	}

	actFile, err := act2.Load(entry.Data)
	if err != nil {
		log.Fatal(err)
	}

	return actFile
}

func BuildSpriteFilePath(gender GenderType, jobSpriteID jobspriteid2.Type) string {
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
