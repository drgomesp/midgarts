// +build glfw

package main

import (
	"fmt"
	"os"

	"github.com/faiface/pixel/pixelgl"
	"github.com/inkyblackness/imgui-go/v3"
	"github.com/project-midgard/midgarts/internal/grfexplorer"
	"github.com/project-midgard/midgarts/internal/platforms"
	"github.com/project-midgard/midgarts/internal/renderers"
)

func main() {
	pixelgl.Run(run)
}

func run() {
	context := imgui.CreateContext(nil)
	defer context.Destroy()
	io := imgui.CurrentIO()

	platform, err := platforms.NewGLFW(io, platforms.GLFWClientAPIOpenGL3)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer platform.Dispose()

	renderer, err := renderers.NewOpenGL3(io)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(-1)
	}
	defer renderer.Dispose()

	grfexplorer.Run(platform, renderer)
}
