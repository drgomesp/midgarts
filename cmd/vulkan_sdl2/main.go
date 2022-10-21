package main

import (
	"os"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
	vk "github.com/vulkan-go/vulkan"
	"github.com/xlab/closer"

	"github.com/project-midgard/midgarts/pkg/core"
)

func init() {

	zerolog.SetGlobalLevel(zerolog.TraceLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	runtime.LockOSThread()
}

func main() {
	checkErr(sdl.Init(sdl.INIT_VIDEO | sdl.INIT_EVENTS))
	defer sdl.Quit()

	checkErr(sdl.VulkanLoadLibrary(""))
	defer sdl.VulkanUnloadLibrary()

	vk.SetGetInstanceProcAddr(sdl.VulkanGetVkGetInstanceProcAddr())
	checkErr(vk.Init())
	defer closer.Close()

	app := core.NewApplication(640, 480)
	reqDim := app.VulkanSwapchainDimensions()
	window, err := sdl.CreateWindow("Vulkan (SDL2)",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(reqDim.Width), int32(reqDim.Height),
		sdl.WINDOW_VULKAN)
	checkErr(err)

	exitC := make(chan struct{}, 2)
	doneC := make(chan struct{}, 2)

	fpsDelay := time.Second / 60
	fpsTicker := time.NewTicker(fpsDelay)
	start := time.Now()
	frames := 0

Loop:
	for {
		log.Printf("FPS: %.2f", float64(frames)/time.Now().Sub(start).Seconds())

		select {
		case <-exitC:
			log.Printf("FPS: %.2f", float64(frames)/time.Now().Sub(start).Seconds())

			_ = window.Destroy()
			fpsTicker.Stop()
			doneC <- struct{}{}
			return
		case <-fpsTicker.C:
			frames++
			var event sdl.Event
			for event = sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
				switch t := event.(type) {
				case *sdl.KeyboardEvent:
					if t.Keysym.Sym == sdl.K_ESCAPE {
						exitC <- struct{}{}
						continue Loop
					}
				case *sdl.QuitEvent:
					exitC <- struct{}{}
					continue Loop
				}
			}
		}
	}

}

func checkErr(err error) {
	if err != nil {
		log.Err(err).Send()
	}
}
