package character

import (
	"log"

	"github.com/project-midgard/midgarts/pkg/character/jobspriteid"

	"golang.org/x/text/encoding/charmap"
)

var jobSpriteNameTable = map[jobspriteid.JobSpriteID]string{
	jobspriteid.Novice:    "",
	jobspriteid.Swordsman: "",
}

func init() {
	var (
		dst []byte
		err error
	)

	dst, err = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC3, 0xCA, 0xBA, 0xB8, 0xC0, 0xDA})
	if err != nil {
		log.Fatal(err)
	}

	jobSpriteNameTable[jobspriteid.Novice] = string(dst)

	dst, err = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB0, 0xCB, 0xBB, 0xE7})
	if err != nil {
		log.Fatal(err)
	}

	jobSpriteNameTable[jobspriteid.Swordsman] = string(dst)
}
