package main

import (
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/project-midgard/midgarts/fileformat/grf"
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

		log.Printf("%#v\n", e)

		fileName := os.Args[2]
		parts := strings.Split(fileName, "\\")
		err = ioutil.WriteFile(parts[len(parts)-1], e.Data.Bytes(), 0644)
		if err != nil {
			log.Fatal(err)
		}
	}
}
