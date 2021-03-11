package main

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/project-midgard/midgarts/internal/fileformat/grf"
	"github.com/project-midgard/midgarts/internal/fileformat/spr"
)

func main() {
	f, err := grf.Load(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	entry := os.Args[2]
	entry = `data\sprite\ork_warrior.spr`
	e, err := f.GetEntry(`data\sprite`, entry)

	if err != nil {
		log.Fatal(err)
	}

	sprFile, err := spr.Load(e.Data)
	if err != nil {
		log.Fatal(err)
	}

	img := sprFile.ImageAt(0)

	outputFile, err := os.Create(strings.Trim(fmt.Sprintf("./out/%s.png", entry), `'`))
	if err != nil {
		log.Fatal(err)
	}

	if err = png.Encode(outputFile, img); err != nil {
		log.Fatal(err)
	}

	_ = outputFile.Close()
}
