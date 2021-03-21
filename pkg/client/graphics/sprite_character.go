package graphics

import (
	"time"

	"github.com/project-midgard/midgarts/pkg/common/character"

	"github.com/EngoEngine/engo/common"
	"github.com/EngoEngine/engo/math"
)

const (
	fpsMultiplier = 1.0

	maleFilePathf   = "data/sprite/%s/%s/³²/%s_³²"
	femaleFilePathf = "data/sprite/%s/%s/¿©/%s_¿©"
)

type CharacterSprite struct {
	Gender character.GenderType

	body *Sprite
}

func NewCharacterSprite(gender character.GenderType, bodySprite *Sprite) *CharacterSprite {
	return &CharacterSprite{
		Gender: gender,
		body:   bodySprite,
	}
}

func NewMonsterSprite(bodySprite *Sprite) *CharacterSprite {
	return &CharacterSprite{
		body: bodySprite,
	}
}

func (s *CharacterSprite) GetActionLayerTexture(actIndex int, layerIndex int) *common.Texture {
	var (
		// TODO should act be taken from BodySprite or somewhere else?
		action                = s.body.act.Actions[actIndex]
		frameIndex            int64
		frameCount            = len(action.Frames)
		timeNeededForOneFrame = float32(action.Delay.Milliseconds()) * 1.0 / fpsMultiplier
	)

	timeAnimationStarted := time.Now()
	timeNeededForOneFrame = math.Max(timeNeededForOneFrame, 100)
	elapsedTime := time.Since(timeAnimationStarted)
	frameIndex = elapsedTime.Milliseconds() / int64(timeNeededForOneFrame)
	frameIndex = frameIndex % int64(frameCount)

	frame := action.Frames[frameIndex]

	// TODO draw all layers?
	layer := frame.Layers[layerIndex]

	return s.body.GetTextureAtIndex(layer.SpriteFrameIndex)
}
