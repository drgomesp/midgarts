package main

import (
	"context"
	"fmt"
	"image/png"
	"log"
	"os"

	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"

	"github.com/EngoEngine/engo"
	"github.com/RaniSputnik/lovepac/packer"
	"github.com/RaniSputnik/lovepac/target"
	"github.com/project-midgard/midgarts/pkg/client/graphics"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
)

var grfFile *grf.File

func main() {
	opts := engo.RunOptions{
		Title:  "Hello World",
		Width:  400,
		Height: 400,
	}
	engo.Run(opts, &myScene{})
}

type myScene struct{}

// Type uniquely defines your game type
func (*myScene) Type() string { return "myGame" }

func (s *myScene) Preload() {
	var err error

	if grfFile, err = grf.Load("/home/drgomesp/grf/data.grf"); err != nil {
		log.Fatal(err)
	}

	jid := jobspriteid.Magician
	var sprite *graphics.Sprite
	//if sprite, err = graphics.LoadCharacterSprite(grfFile, character.Male, jid); err != nil {
	//	log.Fatal(err)
	//}
	if sprite, err = graphics.LoadSprite(grfFile, "data/sprite/ork_warrior"); err != nil {
		log.Fatal(err)
	}

	frames := sprite.SpriteFile.Frames

	var maxWidth, maxHeight uint16
	for _, frame := range frames {
		if frame.Width > maxWidth {
			maxWidth = frame.Width
		}

		if frame.Height > maxHeight {
			maxHeight = frame.Height
		}
	}

	spriteSheetWidth := int(maxWidth) * len(frames) / 5
	spriteSheetHeight := int(maxHeight) * len(frames) / 5
	//
	//var genderStr string
	//if sprite.Gender == character.Male {
	//	genderStr = "m"
	//} else {
	//	genderStr = "f"
	//}

	for i := range sprite.SpriteFile.Frames {
		img := sprite.SpriteFile.ImageAt(i)
		path := fmt.Sprintf("out/ork_warrior/%d.png", i)
		outputFile, err := os.Create(fmt.Sprintf("assets/%s", path))
		if err != nil {
			log.Fatal(err)
		}

		// Encode takes a writer interface and an image interface
		// We pass it the File and the RGBA
		if err = png.Encode(outputFile, img); err != nil {
			log.Fatal(err)
		}
	}

	params := packer.Params{
		Name:       fmt.Sprintf("%d", jid),
		Format:     target.Starling,
		Input:      packer.NewFileStream(fmt.Sprintf("./assets/out/%s", "ork_warrior")),
		Output:     packer.NewFileOutputter(fmt.Sprintf("./assets/build/%s", "ork_warrior")),
		Width:      spriteSheetWidth,
		Height:     spriteSheetHeight,
		MaxAtlases: 1,
	}

	log.Fatal(packer.Run(context.Background(), &params))
}

func (*myScene) Setup(u engo.Updater) {
}
