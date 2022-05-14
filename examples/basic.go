package main

import (
	"image"
	"image/draw"
	"image/jpeg"
	"log"
	"os"
)

func main() {
	f1, err := os.Open("examples/red.jpg")
	checkErr(err)
	defer f1.Close()

	f2, err := os.Open("examples/green.jpg")
	checkErr(err)
	defer f2.Close()

	img1, err := jpeg.Decode(f1)
	checkErr(err)
	sr1 := img1.Bounds()

	var dp image.Point
	img1Rect := image.Rectangle{Min: dp, Max: dp.Add(sr1.Size())}
	_ = img1Rect
	img2, err := jpeg.Decode(f2)
	checkErr(err)
	//sr2 := img2.Bounds()

	//var dp image.Point
	//img2Rect := image.Rectangle{Min: dp, Max: dp.Add(sr2.Size())}
	//_ = img2Rect

	canvasRect := image.Rect(0, 0, 500, 500)
	canvas := image.NewRGBA(canvasRect)

	offset := [2]int{0, 0}

	draw.Draw(canvas, img1.Bounds(), img1, image.Point{}, draw.Over)
	draw.Draw(canvas, img2.Bounds(), img2, image.Point{offset[0], offset[1]}, draw.Over)

	out, err := os.Create("out.jpeg")
	checkErr(err)
	defer out.Close()

	err = jpeg.Encode(out, canvas, nil)
	checkErr(err)
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
