package main

import (
	"fmt"
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

	entry := os.Args[2]
	//entry = `data\sprite\npc\4_f_kafra1.spr`
	e, err := f.GetEntry(entry)

	if err != nil {
		log.Fatal(err)
	}

	sprFile, err := spr.Load(e.Data)
	if err != nil {
		log.Fatal(err)
	}

	img := sprFile.ImageAt(0)

	outputFile, err := os.Create(fmt.Sprintf("./out/%s", entry))
	if err != nil {
		log.Fatal(err)
	}

	if err = png.Encode(outputFile, img); err != nil {
		log.Fatal(err)
	}

	defer func() {
		_ = outputFile.Close()
	}()
}
