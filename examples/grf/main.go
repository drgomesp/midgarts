package main

import (
	"image/png"
	"log"
	"os"

	"github.com/project-midgard/midgarts/pkg/fileformat/spr"

	"github.com/project-midgard/midgarts/pkg/fileformat/grf"
)

func main() {
	f, err := grf.Load(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	entry := os.Args[2]
	e, err := f.GetEntry(entry)

	if err != nil {
		log.Fatal(err)
	}

	sprFile, err := spr.Load(e.Data)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", sprFile)

	img := sprFile.ImageAt(0)

	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("out/test.png")
	if err != nil {
		log.Fatal(err)
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	if err = png.Encode(outputFile, img); err != nil {
		log.Fatal(err)
	}
}
