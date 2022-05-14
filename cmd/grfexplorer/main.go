package main

import (
	g "github.com/AllenDang/giu"
	"github.com/project-midgard/midgarts/internal/character"
	"github.com/project-midgard/midgarts/internal/character/jobspriteid"
	"github.com/project-midgard/midgarts/internal/fileformat/grf"
	"github.com/project-midgard/midgarts/pkg/char"
	"github.com/rs/zerolog/log"
)

var GRF *grf.File

var WidgetChar g.Widget
var SpriteLoader *char.SpriteLoader

func init() {
	var err error
	GRF, err = grf.Load("./assets/grf/data.grf")
	noErr(err)

	SpriteLoader = char.NewSpriteLoader(GRF)
}

func main() {
	wnd := g.NewMasterWindow("Hello world", 640, 480, 0, nil)
	wnd.Run(Run)
}

func Run() {
	go func() {
		sprite, err := SpriteLoader.LoadSprite(character.Female, jobspriteid.Blacksmith, 23, 0)

		if noErr(err) {
			var btex *g.Texture
			btex, err = g.NewTextureFromRgba(sprite.Image)
			noErr(err)

			bsize := sprite.Image.Rect.Size()
			WidgetChar = g.Image(btex).Size(float32(bsize.X), float32(bsize.Y))
		}
	}()

	g.SingleWindowWithMenuBar("splitter").Layout(
		g.MenuBar().Layout(
			g.Menu("File").Layout(
				g.MenuItem("Open"),
				g.MenuItem("Save"),
				// You could add any kind of widget here, not just menu item.
				g.Menu("Save as ...").Layout(
					g.MenuItem("Excel file"),
					g.MenuItem("CSV file"),
					g.Button("Button inside menu"),
				),
			),
		),
		WidgetChar,
	)
}

func noErr(err error) bool {
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	return true
}
