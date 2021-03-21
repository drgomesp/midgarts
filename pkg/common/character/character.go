package character

import (
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"golang.org/x/text/encoding/charmap"
)

var JobSpriteNameTable = map[jobspriteid.Type]string{
	jobspriteid.Novice: "",
	jobspriteid.MonkH:  "",
}

func init() {
	var dst []byte

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC3, 0xCA, 0xBA, 0xB8, 0xC0, 0xDA})
	JobSpriteNameTable[jobspriteid.Novice] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB0, 0xCB, 0xBB, 0xE7})
	JobSpriteNameTable[jobspriteid.Swordsman] = string(dst)

	// TODO check and see if this string actually works
	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte("¸¶¹Ý»Ç"))
	JobSpriteNameTable[jobspriteid.Magician] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB1, 0xC3, 0xBC, 0xF6})
	JobSpriteNameTable[jobspriteid.Archer] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xBC, 0xBA, 0xC1, 0xF7, 0xC0, 0xDA})
	JobSpriteNameTable[jobspriteid.Alcolyte] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xBB, 0xF3, 0xC0, 0xCE})
	JobSpriteNameTable[jobspriteid.Merchant] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB5, 0xB5, 0xB5, 0xCF})
	JobSpriteNameTable[jobspriteid.Thief] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC3, 0xA8, 0xC7, 0xC7, 0xBF, 0xC2})
	JobSpriteNameTable[jobspriteid.MonkH] = string(dst)
}
