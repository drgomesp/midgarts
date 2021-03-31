package main

import (
	"log"
)

func main() {
	var (
		c   *MidgartsClient
		err error
	)

	if c, err = NewMidgartsClient(WithTargetFPS(60)); err != nil {
		log.Fatal(err)
	}

	c.Run()
}
