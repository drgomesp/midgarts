package libvulkan

import (
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
	vk "github.com/vulkan-go/vulkan"
)

type Application struct {
	config   *Config
	window   *sdl.Window
	instance vk.Instance

	name  string
	debug bool
}

func NewApplication(config *Config) (*Application, error) {
	app := &Application{
		config: config,
	}

	if err := app.loadExtensions(); err != nil {
		return nil, err
	}

	return app, nil
}

func (app *Application) loadExtensions() error {
	extensions := app.window.VulkanGetInstanceExtensions()
	if app.debug {
		extensions = append(extensions, "VK_EXT_debug_report")
	}

	log.Info().Strs("available_extensions", extensions).
		Msgf("vulkan available extensions")

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
		log.Info().Msgf("missing required instance extensions during init")
	}

	log.Info().Msgf("enabled vulkan instance extensions")

	// Create instance
	var instance vk.Instance
	ret = vk.CreateInstance(&vk.InstanceCreateInfo{
		SType: vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo: &vk.ApplicationInfo{
			SType:              vk.StructureTypeApplicationInfo,
			ApiVersion:         uint32(app.config.APIVersion),
			ApplicationVersion: uint32(app.config.AppVersion),
			PApplicationName:   safeString(app.name),
			PEngineName:        "libvulkan\x00",
		},
		EnabledExtensionCount:   uint32(len(instanceExtensions)),
		PpEnabledExtensionNames: instanceExtensions,
		//EnabledLayerCount:       uint32(len(validationLayers)),
		//PpEnabledLayerNames:     validationLayers,
	}, nil, &instance)
	if err := VkError(ret); err != nil {
		return err
	}

	app.instance = instance

	if err := vk.InitInstance(instance); err != nil {
		return err
	}

	return nil
}
