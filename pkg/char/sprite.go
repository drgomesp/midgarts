package char

import (
	"github.com/project-midgard/midgarts/internal/character"
	"github.com/project-midgard/midgarts/internal/character/jobspriteid"
	"github.com/project-midgard/midgarts/internal/component"
	"github.com/project-midgard/midgarts/internal/fileformat/grf"
	"image"
	"image/draw"
)

type Sprite struct {
	Image *image.RGBA
}

//
//	jobspriteid.Type
//  direction?
//
//
//

type SpriteLoader struct {
	*grf.File
}

func NewSpriteLoader(grfFile *grf.File) *SpriteLoader {
	return &SpriteLoader{grfFile}
}

func (s *SpriteLoader) LoadSprite(
	gender character.GenderType,
	jid jobspriteid.Type,
	headIndex character.HeadIndex,
	spriteIndex character.SpriteIndex,
) (*Sprite, error) {
	attachments, err := component.NewCharacterAttachmentComponent(
		s.File,
		component.CharacterAttachmentComponentConfig{
			Gender:      gender,
			JobSpriteID: jid,
			HeadIndex:   headIndex,
		},
	)
	if err != nil {
		return nil, err
	}

	var output, prevImg *image.RGBA
	for _, attachment := range attachments.Files {
		if attachment.SPR != nil {
			if img := attachment.SPR.ImageAt(spriteIndex); img != nil {
				img := img.RGBA

				if prevImg == nil {
					output := image.NewRGBA(image.Rect(0, 0, 300, 300))
					draw.Draw(output, img.Bounds(), img, image.Point{}, draw.Src)
				} else {
					point := image.Point{Y: prevImg.Bounds().Dy()}
					rect := image.Rectangle{Min: point, Max: point.Add(img.Bounds().Size())}
					targetRect := image.Rectangle{Min: image.Point{}, Max: rect.Max}

					if output == nil {
						output = image.NewRGBA(targetRect)
					}

					draw.Draw(output, rect, img, image.Point{}, draw.Src)
				}

				prevImg = img
			}
		}
	}

	return &Sprite{output}, nil
}
