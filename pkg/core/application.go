package core

import (
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
	as "github.com/vulkan-go/asche"
	vk "github.com/vulkan-go/vulkan"
)

type Application struct {
	debug               bool
	winWidth, winHeight uint32
	window              *sdl.Window
}

func NewApplication(winWidth, winHeight uint32) *Application {
	app := &Application{
		winWidth:  winWidth,
		winHeight: winHeight,
	}

	return app
}

func (a *Application) VulkanSurface(instance vk.Instance) (surface vk.Surface) {
	surfPtr, err := a.window.VulkanCreateSurface(instance)
	if err != nil {
		log.Error().Err(err).Msgf("vulkan error")
		return vk.NullSurface
	}

	surf := vk.SurfaceFromPointer(uintptr(surfPtr))
	return surf
}

func (a *Application) VulkanDebug() bool {
	return a.debug
}

func (a *Application) VulkanDeviceExtensions() []string {
	return []string{
		"VK_KHR_swapchain",
	}
}

func (a *Application) VulkanSwapchainDimensions() *as.SwapchainDimensions {
	return &as.SwapchainDimensions{
		Width: a.winWidth, Height: a.winHeight, Format: vk.FormatB8g8r8a8Unorm,
	}
}

func (a *Application) VulkanInstanceExtensions() []string {
	extensions := a.window.VulkanGetInstanceExtensions()
	if a.debug {
		extensions = append(extensions, "VK_EXT_debug_report")
	}

	return extensions
}
