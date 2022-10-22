package libvulkan

import (
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
	vk "github.com/vulkan-go/vulkan"
)

type Application struct {
	window *sdl.Window
	debug  bool
}

func NewApplication() *Application {
	app := &Application{}

	app.loadExtensions()

	return app
}

func (a *Application) loadExtensions() error {
	extensions := a.window.VulkanGetInstanceExtensions()
	if a.debug {
		extensions = append(extensions, "VK_EXT_debug_report")
	}

	log.Info().Strs("available_extensions", extensions).Msgf("vulkan: loaded extension info")

	requiredInstanceExtensions := safeStrings(extensions)
	var actualInstanceExtensions []string

	var count uint32
	ret := vk.EnumerateInstanceExtensionProperties("", &count, nil)
	if err := VkError(ret); err != nil {
		return err
	}

	list := make([]vk.ExtensionProperties, count)
	ret = vk.EnumerateInstanceExtensionProperties("", &count, list)
	if err := VkError(ret); err != nil {
		return err
	}

	for _, ext := range list {
		ext.Deref()
		actualInstanceExtensions = append(actualInstanceExtensions, vk.ToString(ext.ExtensionName[:]))
	}
	if err := VkError(ret); err != nil {
		return err
	}

	instanceExtensions, missing := checkExisting(actualInstanceExtensions, requiredInstanceExtensions)
	if missing > 0 {
		log.Info().Msgf("vulkan warning: missing", missing, "required instance extensions during init")
	}

	log.Info().Msgf("vulkan: enabling %d instance extensions", len(instanceExtensions))

	return nil
}
