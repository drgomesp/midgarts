package main

import (
	"log"
	"os"

	"github.com/project-midgard/midgarts/internal/fileformat/act"
	"github.com/project-midgard/midgarts/internal/fileformat/grf"
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

	actFile, err := act.Load(e.Data)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%#v\n", actFile)
}
