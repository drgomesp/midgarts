package main

import (
	"os"
	"runtime"

	"github.com/davecgh/go-spew/spew"
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

	app := libvulkan.NewApplication()
	spew.Dump(app)
}

func checkErr(err error) {
	if err != nil {
		log.Err(err).Send()
	}
}
