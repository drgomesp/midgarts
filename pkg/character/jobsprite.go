package character

import (
	"github.com/project-midgard/midgarts/pkg/character/jobspriteid"
	"golang.org/x/text/encoding/charmap"
)

var JobSpriteNameTable = map[jobspriteid.Type]string{}

func init() {
	var dst []byte

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC3, 0xCA, 0xBA, 0xB8, 0xC0, 0xDA})
	JobSpriteNameTable[jobspriteid.Novice] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB0, 0xCB, 0xBB, 0xE7})
	JobSpriteNameTable[jobspriteid.Swordsman] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB8, 0xB6, 0xB9, 0xDD, 0xBB, 0xC7})
	JobSpriteNameTable[jobspriteid.Magician] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB1, 0xC3, 0xBC, 0xF6})
	JobSpriteNameTable[jobspriteid.Archer] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xBC, 0xBA, 0xC1, 0xF7, 0xC0, 0xDA})
	JobSpriteNameTable[jobspriteid.Alcolyte] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xBB, 0xF3, 0xC0, 0xCE})
	JobSpriteNameTable[jobspriteid.Merchant] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB5, 0xB5, 0xB5, 0xCF})
	JobSpriteNameTable[jobspriteid.Thief] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB8, 0xF9, 0xC5, 0xA9})
	JobSpriteNameTable[jobspriteid.Monk] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB1, 0xE2, 0xBB, 0xE7})
	JobSpriteNameTable[jobspriteid.Knight] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC6, 0xE4, 0xC4, 0xDA, 0xC6, 0xE4, 0xC4, 0xDA, 0x5f, 0xB1, 0xE2, 0xBB, 0xE7})
	JobSpriteNameTable[jobspriteid.Knight2] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC7, 0xC1, 0xB8, 0xAE, 0xBD, 0xBA, 0xC6, 0xAE})
	JobSpriteNameTable[jobspriteid.Priest] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC0, 0xA7, 0xC0, 0xFA, 0xB5, 0xE5})
	JobSpriteNameTable[jobspriteid.Wizard] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC1, 0xA6, 0xC3, 0xB6, 0xB0, 0xF8})
	JobSpriteNameTable[jobspriteid.Blacksmith] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC7, 0xE5, 0xC5, 0xCD})
	JobSpriteNameTable[jobspriteid.Hunter] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC5, 0xA9, 0xB7, 0xE7, 0xBC, 0xBC, 0xC0, 0xCC, 0xB4, 0xF5})
	JobSpriteNameTable[jobspriteid.Crusader] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xBD, 0xC5, 0xC6, 0xE4, 0xC4, 0xDA, 0xC5, 0xA9, 0xB7, 0xE7, 0xBC, 0xBC, 0xC0, 0xCC, 0xB4, 0xF5})
	JobSpriteNameTable[jobspriteid.Crusader2] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xBC, 0xBC, 0xC0, 0xCC, 0xC1, 0xF6})
	JobSpriteNameTable[jobspriteid.Sage] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB7, 0xCE, 0xB1, 0xD7})
	JobSpriteNameTable[jobspriteid.Rogue] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xBF, 0xAC, 0xB1, 0xDD, 0xBC, 0xFA, 0xBB, 0xE7})
	JobSpriteNameTable[jobspriteid.Alchemist] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xBE, 0xEE, 0xBC, 0xBC, 0xBD, 0xC5})
	JobSpriteNameTable[jobspriteid.Assassin] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB9, 0xD9, 0xB5, 0xE5})
	JobSpriteNameTable[jobspriteid.Bard] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xB9, 0xAB, 0xC8, 0xF1})
	JobSpriteNameTable[jobspriteid.Dancer] = string(dst)

	dst, _ = charmap.Windows1252.NewDecoder().Bytes([]byte{0xC3, 0xA8, 0xC7, 0xC7, 0xBF, 0xC2})
	JobSpriteNameTable[jobspriteid.MonkH] = string(dst)
}
