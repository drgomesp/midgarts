package char

//
//import (
//	"image"
//	"image/color"
//
//	"github.com/project-midgard/midgarts/internal/character"
//)
//
//type Accessory struct {
//	Color         *image.Uniform
//	Type          character.AttachmentType
//	Offset        image.Point
//	PositionFrame image.Point
//	PositionLayer image.Point
//	Image         *image.RGBA
//}
//
//func NewAccessory(
//	elem character.AttachmentType,
//	offset,
//	positionFrame,
//	positionLayer image.Point,
//	img *image.RGBA,
//) Accessory {
//	var c *image.Uniform
//
//	Blue := color.RGBA{R: 0, G: 0, B: 255, A: 255}
//	Green := color.RGBA{R: 0, G: 255, B: 0, A: 255}
//	Red := color.RGBA{R: 255, G: 0, B: 0, A: 255}
//
//	switch elem {
//	case character.AttachmentShadow:
//		c = &image.Uniform{C: Red}
//	case character.AttachmentBody:
//		c = &image.Uniform{C: Green}
//	case character.AttachmentHead:
//		c = &image.Uniform{C: Blue}
//	}
//
//	return Accessory{
//		c,
//		elem,
//		offset,
//		positionFrame,
//		positionLayer,
//		img,
//	}
//}
//
//// Anchor takes a set of accessories, calculates the proper anchor
//// points for each one of them and returns a new set of accessories
//// with the calculated anchor points.
//func Anchor(accessories ...Accessory) []Accessory {
//	var anchor image.Point
//
//	res := make([]Accessory, 0)
//	for _, acc := range accessories {
//		var pos image.Point
//
//		if acc.Type != character.AttachmentBody &&
//			acc.Type != character.AttachmentShield {
//			pos.X = anchor.X - acc.PositionFrame.X
//			pos.Y = anchor.Y - acc.PositionFrame.Y
//		}
//
//		acc.Offset.X = acc.PositionFrame.X + pos.X
//		acc.Offset.Y = acc.PositionFrame.Y + pos.Y
//
//		res = append(res, acc)
//
//		anchor.X = acc.PositionFrame.X
//		anchor.Y = acc.PositionFrame.Y
//	}
//
//	return res
//}
