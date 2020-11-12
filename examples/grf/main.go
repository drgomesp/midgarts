package main

import (
	"image"
	"image/color"
	"image/draw"
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

	if len(os.Args) > 2 {
		e, err := f.GetEntry(os.Args[2])
		if err != nil {
			log.Fatal(err)
		}

		sprFile, err := spr.Load(e.Data)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("entry=%#v\n", sprFile)

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

		//imgPath := "foo.png"
		//out, err := os.Create(imgPath)
		//if err != nil {
		//	log.Fatal(err)
		//}
		//w, h := int(sprFile.Frames[0].Width), int(sprFile.Frames[0].Height)
		//img := createImage(w, h)
		//data := sprFile.Frames[0].Data
		//outputWidth := sprFile.Frames[0].Width
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
		i := 0
		for fileName, _ := range f.GetEntries() {
			e, err := f.GetEntry(fileName)
			if err != nil {
				log.Println(err)
				continue
			}

			sprFile, err := spr.Load(e.Data)

			if err != nil {
				log.Println(err)
				continue
			}

			log.Printf("[%s] %#v\n", fileName, sprFile)

			if i > 5 {
				break
			}

			i++

			log.Print("\n\n")
		}
	}

}

func createImage(width int, height int) *image.RGBA {
	rect := image.Rect(0, 0, width, height)
	img := image.NewRGBA(rect)
	background := color.RGBA{0, 0xFF, 0, 0xCC}
	draw.Draw(img, img.Bounds(), &image.Uniform{background}, image.Point{}, draw.Src)
	return img
}
