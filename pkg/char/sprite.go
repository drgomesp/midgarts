package char

import (
	"github.com/project-midgard/midgarts/internal/character"
	"github.com/project-midgard/midgarts/internal/character/jobspriteid"
	"github.com/project-midgard/midgarts/internal/component"
	"github.com/project-midgard/midgarts/internal/fileformat/grf"
	"github.com/rs/zerolog/log"
	"image"
	"image/draw"
)

type Sprite struct {
	Image *image.RGBA
}

type SpriteOptions struct {
	Gender      character.GenderType
	JobSpriteID jobspriteid.Type
	HeadIndex   character.HeadIndex
	SpriteIndex character.SpriteIndex
}

type SpriteLoader struct {
	*grf.File
}

func NewSpriteLoader(grfFile *grf.File) *SpriteLoader {
	return &SpriteLoader{grfFile}
}

func (s *SpriteLoader) LoadSprite(opts SpriteOptions) (*Sprite, error) {
	attachments, err := component.NewCharacterAttachmentComponent(
		s.File,
		component.CharacterAttachmentComponentConfig{
			Gender:      opts.Gender,
			JobSpriteID: opts.JobSpriteID,
			HeadIndex:   opts.HeadIndex,
		},
	)
	if err != nil {
		return nil, err
	}

	var offset, position [2]float32
	var canvas, prevImg *image.RGBA
	for t, attachment := range attachments.Files {
		if attachment.SPR != nil {
			elem := character.AttachmentType(t)

			action := attachment.ACT.Actions[0]
			frame := action.Frames[0]
			layer := frame.Layers[0]

			if frame = action.Frames[0]; len(frame.Layers) == 0 {
				offset = [2]float32{0, 0}
				continue
			}

			if len(frame.Positions) > 0 &&
				elem != character.AttachmentBody &&
				elem != character.AttachmentShield {
				position[0] = offset[0] - float32(frame.Positions[0][0])
				position[1] = offset[1] - float32(frame.Positions[0][1])
			}

			if img := attachment.SPR.ImageAt(opts.SpriteIndex); img != nil {
				img := img.RGBA

				pos := [2]int{
					int(position[0]) + int(float32(layer.Position[0])+offset[0]),
					int(position[0]) + int(float32(layer.Position[1])+offset[1]),
				}

				if prevImg == nil {
					canvas = image.NewRGBA(image.Rect(pos[0], pos[1], 300, 300))
					draw.Draw(canvas, img.Bounds(), img, image.Point{}, draw.Over)
				} else {
					log.Trace().Msgf("attachment[%s]: pos=%v", elem, position)
					targetRect := image.Rect(pos[0], pos[1], 100, 100)

					_ = pos
					draw.Draw(canvas, targetRect, img, image.Point{pos[0], pos[1]}, draw.Over)

					// Save offset reference
					if len(frame.Positions) > 0 {
						offset = [2]float32{
							float32(frame.Positions[0][0]),
							float32(frame.Positions[0][1]),
						}
					}

				}

				prevImg = img
			}
		}
	}

	return &Sprite{canvas}, nil
}
