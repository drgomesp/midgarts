package main

import (
	"log"
)

func main() {
	var (
		c   *MidgartsClient
		err error
	)

	if c, err = NewMidgartsClient(WithTargetFPS(3)); err != nil {
		log.Fatal(err)
	}

	c.Run()
}
