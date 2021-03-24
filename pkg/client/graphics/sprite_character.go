package graphics

import (
	"github.com/project-midgard/midgarts/pkg/common/character"
)

type CharacterSprite struct {
	Gender character.GenderType

	Body *Sprite
}

func NewCharacterSprite(gender character.GenderType, bodySprite *Sprite) *CharacterSprite {
	return &CharacterSprite{
		Gender: gender,
		Body:   bodySprite,
	}
}
