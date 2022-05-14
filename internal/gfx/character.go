package gfx

//
//import (
//	character2 "github.com/project-midgard/midgarts/internal/character"
//	"github.com/project-midgard/midgarts/internal/character/jobspriteid"
//	"github.com/project-midgard/midgarts/internal/component"
//	"github.com/project-midgard/midgarts/internal/fileformat/grf"
//	"image"
//	"image/draw"
//)
//
//// CharSpriteImage returns an image of a character sprite in a given index (action/frame).
//func CharSpriteImage(
//	grfFile *grf.File,
//	gender character2.GenderType,
//	jobSpriteID jobspriteid.Type,
//	index int,
//) (*image.RGBA, error) {
//	attachments, err := component.NewCharacterAttachmentComponent(grfFile, component.CharacterAttachmentComponentConfig{
//		Gender:           gender,
//		JobSpriteID:      jobSpriteID,
//		HeadIndex:        23,
//		EnableShield:     false,
//		ShieldSpriteName: "",
//	})
//	if err != nil {
//		return nil, err
//	}
//
//	head := attachments.Files[character2.AttachmentHead]
//	img1 := head.SPR.ImageAt(index).RGBA
//
//	body := attachments.Files[character2.AttachmentBody]
//	img2 := body.SPR.ImageAt(index).RGBA
//	sp2 := image.Point{0, img1.Bounds().Dy()}
//	r2 := image.Rectangle{sp2, sp2.Add(img2.Bounds().Size())}
//
//	r := image.Rectangle{
//		Min: image.Point{0, 0},
//		Max: r2.Max,
//	}
//
//	result := image.NewRGBA(r)
//
//	draw.Draw(result, img1.Bounds(), img1, image.Point{0, 0}, draw.Src)
//	draw.Draw(result, r2, img2, image.Point{0, 0}, draw.Src)
//
//	return result, nil
//}
