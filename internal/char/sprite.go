package char

//
//import (
//	"image"
//	"image/draw"
//
//	"github.com/rs/zerolog/log"
//
//	"github.com/project-midgard/midgarts/internal/character"
//	"github.com/project-midgard/midgarts/internal/character/jobspriteid"
//	"github.com/project-midgard/midgarts/internal/component"
//	"github.com/project-midgard/midgarts/internal/fileformat/grf"
//)
//
//type Sprite struct {
//	Image *image.RGBA
//}
//
//type SpriteLoader struct {
//	*grf.File
//}
//
//func NewSpriteLoader(grfFile *grf.File) *SpriteLoader {
//	return &SpriteLoader{grfFile}
//}
//
//func (s *SpriteLoader) LoadSprite(
//	gender character.GenderType,
//	jid jobspriteid.Type,
//	headIndex character.HeadIndex,
//	spriteIndex character.SpriteIndex,
//) (*Sprite, error) {
//	attachments, err := component.NewCharacterAttachmentComponent(
//		s.File,
//		component.CharacterAttachmentComponentConfig{
//			Gender:      gender,
//			JobSpriteID: jid,
//			HeadIndex:   headIndex,
//		},
//	)
//	if err != nil {
//		return nil, err
//	}
//
//	// Set the shader to the first
//	// Set action index to the default "idle" action
//	actionIndex := 0
//	// Set the frame index to the first action frame
//	frameIndex := 0
//	// Set the base layer
//	layerIndex := 0
//
//	accessories := make([]Accessory, character.NumAttachments-1)
//	for t, att := range attachments.Files {
//		action := att.ACT.Actions[actionIndex]
//		frame := action.Frames[frameIndex]
//
//		framePos := image.Point{
//			X: int(frame.Positions[0][0]),
//			Y: int(frame.Positions[0][1]),
//		}
//		layerPos := image.Point{
//			X: int(frame.Layers[layerIndex].Position[0]),
//			Y: int(frame.Layers[layerIndex].Position[1]),
//		}
//
//		accessories[t] = NewAccessory(
//			t, image.Point{}, framePos, layerPos, att.SPR.ImageAt(0).RGBA,
//		)
//	}
//
//	accessories = Anchor(accessories...)
//
//	w := 50
//	h := 80
//	canvas := image.NewRGBA(image.Rect(0, 0, w, h))
//
//	draw.Draw(canvas, canvas.Bounds(), &image.Uniform{C: image.White}, image.ZP, draw.Over)
//
//	//var anchor image.Point
//	for _, acc := range accessories {
//		origin := image.Point{}
//		size := acc.Image.Bounds()
//
//		rect := image.Rect(0, canvas.Bounds().Dy()-size.Dy(), size.Dx(), canvas.Bounds().Dy())
//
//		//anchor = acc.Offset
//		draw.Draw(canvas, rect.Add(acc.Offset), acc.Color, origin, draw.Over)
//		//draw.Draw(canvas, rect.Add(acc.Offset), acc.Image, origin, draw.Over)
//		//draw.Draw(canvas, rect, acc.Image, origin.Sub(offset), draw.Over)
//		//size := acc.Image.Rect.Size()
//		//origin := image.Point{X: 0, Y: canvas.Bounds().Dy() - size.Y}
//		//
//		////var origin image.Point
//		//
//		////pos := image.Point{X: origin.X + acc.PositionAnchor.X, Y: origin.Y + acc.PositionAnchor.Y}
//		////origin = origin.Add(acc.PositionAnchor.Point())
//		////origin = image.Point{X: origin.X + acc.PositionAnchor.X, Y: origin.Y + acc.PositionAnchor.Y}
//		//
//
//		//
//		//rect := image.Rect(origin.X, origin.Y, size.X, size.Y)
//		////offset := image.Point{}
//		//draw.Draw(canvas, rect, acc.Color, origin, draw.Over)
//		////draw.Draw(canvas, canvas.Rect, acc.Image, image.ZP, draw.Over)
//	}
//
//	log.Debug().Msg("\n")
//
//	return &Sprite{canvas}, nil
//}
