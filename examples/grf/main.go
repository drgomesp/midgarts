package main

import (
	"image/png"
	"log"
	"os"

	"github.com/project-midgard/midgarts/fileformat/grf"
	"github.com/project-midgard/midgarts/fileformat/spr"
)

func main() {
	f, err := grf.NewFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	e, err := f.GetEntry(`data\sprite\ork_warrior.spr`)
	//e, err := f.GetEntry(`data\sprite\count.spr`)
	//e, err := f.GetEntry(`data\sprite\npc\bigfoot.spr`)
	if err != nil {
		log.Fatal(err)
	}

	sprFile, err := spr.Load(e.Data)
	if err != nil {
		log.Fatal(err)
	}

	img := sprFile.ImageAt(0)

	outputFile, err := os.Create("./out/test.png")
	if err != nil {
		log.Fatal(err)
	}

	if err = png.Encode(outputFile, img); err != nil {
		log.Fatal(err)
	}

	outputFile.Close()
}
