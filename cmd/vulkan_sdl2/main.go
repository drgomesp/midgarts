package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
	vk "github.com/vulkan-go/vulkan"
	"github.com/xlab/closer"

	"github.com/project-midgard/midgarts/pkg/libvulkan"
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

	window, err := sdl.CreateWindow("libvulkan (SDL2)",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		800, 600,
		sdl.WINDOW_VULKAN)
	checkErr(err)

	app, err := libvulkan.NewApplication(libvulkan.DefaultConfig(), window)
	checkErr(err)

	// some sync logic
	doneC := make(chan struct{}, 2)
	exitC := make(chan struct{}, 2)
	defer closer.Bind(func() {
		exitC <- struct{}{}
		<-doneC
		log.Info().Msg("Bye!")
	})

	fpsDelay := time.Second / 60
	fpsTicker := time.NewTicker(fpsDelay)
	start := time.Now()
	frames := 0
_MainLoop:
	for {
		select {
		case <-exitC:
			fmt.Printf("FPS: %.2f\n", float64(frames)/time.Now().Sub(start).Seconds())
			//app.Destroy()
			//platform.Destroy()
			//window.Destroy()
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
						continue _MainLoop
					}
				case *sdl.QuitEvent:
					exitC <- struct{}{}
					continue _MainLoop
				}
			}

			_ = app
		}
	}
}

func checkErr(err error) {
	if err != nil {
		log.Err(err).Send()
	}
}