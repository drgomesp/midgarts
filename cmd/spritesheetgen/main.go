package main

import (
	"context"
	"fmt"
	"github.com/EngoEngine/engo"
	"github.com/RaniSputnik/lovepac/packer"
	"github.com/RaniSputnik/lovepac/target"
	"github.com/project-midgard/midgarts/internal/graphics"
	"github.com/project-midgard/midgarts/pkg/common/character"
	"github.com/project-midgard/midgarts/pkg/common/character/jobspriteid"
	"github.com/project-midgard/midgarts/pkg/common/fileformat/grf"
	"image/png"
	"log"
	"os"
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

	jid := jobspriteid.MonkH
	var sprite *graphics.CharacterSprite
	if sprite, err = graphics.LoadCharacterSprite(grfFile, character.Male, jid); err != nil {
		log.Fatal(err)
	}

	frames := sprite.Body.SpriteFile.Frames

	var maxWidth, maxHeight uint16
	for _, frame := range frames {
		if frame.Width > maxWidth {
			maxWidth = frame.Width
		}

		if frame.Height > maxHeight {
			maxHeight = frame.Height
		}
	}

	spriteSheetWidth := int(maxWidth) * len(frames) / 10
	spriteSheetHeight := int(maxHeight) * len(frames) / 10

	var genderStr string
	if sprite.Gender == character.Male {
		genderStr = "m"
	} else {
		genderStr = "f"
	}

	for i, _ := range sprite.Body.SpriteFile.Frames {
		img := sprite.Body.SpriteFile.ImageAt(i)

		fileName := fmt.Sprintf("%d.png", i)
		path := fmt.Sprintf("out/%d/%s/%s", jid, genderStr, fileName)
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
		Input:      packer.NewFileStream(fmt.Sprintf("./assets/out/%d/%s", jid, genderStr)),
		Output:     packer.NewFileOutputter(fmt.Sprintf("./assets/build/%s", genderStr)),
		Width:      spriteSheetWidth,
		Height:     spriteSheetHeight,
		MaxAtlases: 1,
	}

	log.Fatal(packer.Run(context.Background(), &params))
}

func (*myScene) Setup(u engo.Updater) {
}
