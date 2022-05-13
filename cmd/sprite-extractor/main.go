package main

import (
	"github.com/drgomesp/midgarts/pkg/fileformat/grf"
	"github.com/rs/zerolog/log"
	"image/png"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	f, err := grf.Load("./assets/grf/data.grf")
	checkErr(err)

	sprites := make([]*grf.Entry, 0)
	f.GetEntryTree().Traverse(f.GetEntryTree().Root, func(node *grf.EntryTreeNode) {
		for _, entry := range node.Data {
			ext := filepath.Ext(entry.Name)
			//name := strings.TrimSuffix(entry.Name, ext)

			if ext == ".act" || ext == ".spr" {
				if strings.Contains(entry.Name, "viking") {
					panic("OOOOOO")
				}

				e, err := f.GetEntry(entry.Name)
				checkErr(err)
				log.Info().Msgf(entry.Name)

				sprites = append(sprites, e)
			}
		}
	})

	e, err := f.GetEntry("data/sprite/ork_warrior.spr")
	checkErr(err)

	files, err := f.GetActionAndSpriteFiles(strings.TrimSuffix(e.Name, filepath.Ext(e.Name)))
	checkErr(err)

	img := files.SPR.ImageAt(0)

	imgFile, _ := os.Create("out.png")
	err = png.Encode(imgFile, img)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}
