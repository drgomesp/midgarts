package main

import (
	"github.com/EngoEngine/engo"
)

type myScene struct{}

// Type uniquely defines your game type
func (*myScene) Type() string { return "myGame" }

// Preload is called before loading any assets from the disk,
// to allow you to register / queue them
func (*myScene) Preload() {}

// Setup is called before the main loop starts. It allows you
// to add entities and systems to your Scene.
func (*myScene) Setup(engo.Updater) {}

func main() {
	opts := engo.RunOptions{
		Title:  "Hello World",
		Width:  400,
		Height: 400,
	}
	engo.Run(opts, &myScene{})
}
