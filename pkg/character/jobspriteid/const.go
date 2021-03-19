package jobspriteid

import (
	"log"

	"golang.org/x/text/encoding/charmap"
)

type JobSpriteID int

const (
	Novice JobSpriteID = 0
)

var JobSpriteNameTable = map[JobSpriteID]string{}

func init() {
	dst, err := charmap.Windows1252.NewDecoder().Bytes([]byte{0xC3, 0xCA, 0xBA, 0xB8, 0xC0, 0xDA})
	if err != nil {
		log.Fatal(err)
	}

	JobSpriteNameTable[Novice] = string(dst)
}
