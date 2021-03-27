// +build glfw

package main

import (
	g "github.com/AllenDang/giu"
)

func main() {
	wnd := g.NewMasterWindow("Hello world", 640, 480, 0, nil)
	wnd.Run(Run)
}
