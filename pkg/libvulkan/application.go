package libvulkan

import (
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
	vk "github.com/vulkan-go/vulkan"
)

var nilExtensions []string
var nilValidationLayers []string

type Application struct {
	config Config
	window *sdl.Window
	device *Device

	name  string
	debug bool

	extensions       []string
	validationLayers []string
}

func NewApplication(config Config, window *sdl.Window) (*Application, error) {
	app := &Application{
		config: config,
		window: window,
	}

	extensions, err := app.loadExtensions()
	if err != nil {
		return nil, err
	}

	validationLayers, err := app.loadValidationLayers()
	if err != nil {
		return nil, err
	}

	app.device, err = NewDevice(config, extensions, validationLayers, window)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *Application) loadExtensions() ([]string, error) {
	extensions := app.window.VulkanGetInstanceExtensions()
	if app.debug {
		extensions = append(extensions, "VK_EXT_debug_report")
	}

	log.Debug().Strs("instance_extensions", extensions).
		Msgf("instance extensions")

	requiredInstanceExtensions := safeStrings(extensions)
	var actualInstanceExtensions []string

	var count uint32
	ret := vk.EnumerateInstanceExtensionProperties("", &count, nil)
	if err := VkError(ret); err != nil {
		return nilExtensions, err
	}

	list := make([]vk.ExtensionProperties, count)
	ret = vk.EnumerateInstanceExtensionProperties("", &count, list)
	if err := VkError(ret); err != nil {
		return nilExtensions, err
	}

	for _, ext := range list {
		ext.Deref()
		actualInstanceExtensions = append(actualInstanceExtensions, vk.ToString(ext.ExtensionName[:]))
	}
	if err := VkError(ret); err != nil {
		return nilExtensions, err
	}

	instanceExtensions, missing := checkExisting(actualInstanceExtensions, requiredInstanceExtensions)
	if missing > 0 {
		log.Info().Msgf("missing required instance extensions during init")
	}

	log.Info().Msgf("enabled instance extensions")

	return instanceExtensions, nil
}

func (app *Application) loadValidationLayers() ([]string, error) {
	requiredValidationLayers := safeStrings(app.config.ValidationLayers)

	var actualInstanceValidationLayers []string
	var count uint32
	ret := vk.EnumerateInstanceLayerProperties(&count, nil)
	if err := VkError(ret); err != nil {
		return nilValidationLayers, err
	}

	list := make([]vk.LayerProperties, count)
	ret = vk.EnumerateInstanceLayerProperties(&count, list)
	if err := VkError(ret); err != nil {
		return nilValidationLayers, err
	}

	for _, layer := range list {
		layer.Deref()
		actualInstanceValidationLayers = append(actualInstanceValidationLayers, vk.ToString(layer.LayerName[:]))
	}

	validationLayers, missing := checkExisting(actualInstanceValidationLayers, requiredValidationLayers)
	if missing > 0 {
		log.Warn().Msgf("missing %d required validation layers during init", missing)
	}

	log.Debug().Strs("validation_layers", validationLayers).
		Msgf("instance validation layers")

	log.Info().Msgf("enabled instance validation layers")

	return validationLayers, nil
}