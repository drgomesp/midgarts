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

	e, err := f.GetEntry(`data\sprite\npc\1_etc_01.spr`)
	//e, err := f.GetEntry(`data\sprite\count.spr`)
	//e, err := f.GetEntry(`data\sprite\npc\bigfoot.spr`)
	if err != nil {
		log.Fatal(err)
	}

	sprFile, err := spr.Load(e.Data)
	if err != nil {
		log.Fatal(err)
	}

	img := sprFile.Image(0, false, nil)

	// outputFile is a File type which satisfies Writer interface
	outputFile, err := os.Create("test.png")
	if err != nil {
		log.Fatal(err)
	}

	// Encode takes a writer interface and an image interface
	// We pass it the File and the RGBA
	if err = png.Encode(outputFile, img); err != nil {
		log.Fatal(err)
	}

	// Don't forget to close files
	outputFile.Close()

	//fileName := os.Args[2]
	//parts := strings.Split(fileName, "\\")
	//err = ioutil.WriteFile(parts[len(parts)-1], e.Data.Bytes(), 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//fileName := os.Args[2]
	//parts := strings.Split(fileName, "\\")
	//err = ioutil.WriteFile(parts[len(parts)-1], e.Data.Bytes(), 0644)
	//if err != nil {
	//	log.Fatal(err)
	//}

	//
	//for y := 0; y < h; y++ {
	//	for x := 0; x < w; x++ {
	//		idx1 := data[y*w+x] * 4
	//
	//		img.Set(idx1+0, int(data[idx2+0]), color.RGBA{})
	//		img.Set(idx1+1, int(data[idx2+1]), color.RGBA{})
	//		img.Set(idx1+2, int(data[idx2+2]), color.RGBA{})
	//		img.Set(idx1+3, int(data[idx2+3]), color.RGBA{})
	//	}
	//}
	//
	//err = png.Encode(out, img)
}
