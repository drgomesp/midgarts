package graphics

import (
	"github.com/EngoEngine/engo/common"
)

type SpritesheetResource struct {
	Spritesheet *common.Spritesheet
}

func NewSpritesheetResource(spritesheet *common.Spritesheet) *SpritesheetResource {
	return &SpritesheetResource{
		Spritesheet: spritesheet,
	}
}
