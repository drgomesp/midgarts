package main

import (
	"image/png"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/project-midgard/midgarts/fileformat/grf"
	"github.com/project-midgard/midgarts/fileformat/spr"
)

func main() {
	f, err := grf.NewFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	if len(os.Args) > 2 {
		e, err := f.GetEntry(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		sprFile, err := spr.Load(e.Data)
		if err != nil {
			log.Fatal(err)
		}

		sprFile.ToRGBA()
		sprFile.Compile()
		img := sprFile.Image(0, false)

		// outputFile is a File type which satisfies Writer interface
		outputFile, err := os.Create("test.png")
		if err != nil {
			log.Fatal(err)
		}

		// Encode takes a writer interface and an image interface
		// We pass it the File and the RGBA
		png.Encode(outputFile, img)

		// Don't forget to close files
		outputFile.Close()

		spew.Dump(img)

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
	} else {
		//i := 0
		//for fileName, _ := range f.GetEntries() {
		//	e, err := f.GetEntry(fileName)
		//	if err != nil {
		//		continue
		//	}
		//
		//	sprFile, err := spr.Load(e.Data)
		//
		//	if err != nil {
		//		continue
		//	}
		//
		//	//spew.Dump(sprFile)
		//	//log.Printf("[%s] %#v\n", fileName, sprFile)
		//
		//	img := sprFile.Image(0, false, color.RGBA{
		//		R: 155,
		//		G: 10,
		//		B: 15,
		//		A: 0,
		//	})
		//
		//	_ = img
		//	//spew.Dump(img)
		//
		//	i++
		//
		//	if !strings.Contains(fileName, "ork_warrior") {
		//		continue
		//	}
		//
		//	log.Printf("%s\n", fileName)
		//	log.Print("\n\n")
		//}
	}

}
