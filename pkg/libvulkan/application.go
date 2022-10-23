package libvulkan

import (
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
)

var nilExtensions []string
var nilValidationLayers []string

type Application struct {
	config Config
	window *sdl.Window
	device *Device

	name string

	extensions       []string
	validationLayers []string
}

func NewApplication(config Config, window *sdl.Window) (*Application, error) {
	app := &Application{
		config: config,
		window: window,
	}

	if err := app.loadExtensions(); err != nil {
		return nil, err
	}

	err := app.loadValidationLayers()
	if err != nil {
		return nil, err
	}

	app.device, err = NewDevice(config, app.extensions, app.validationLayers, window)
	if err != nil {
		return nil, err
	}

	return app, nil
}

func (app *Application) loadExtensions() error {
	requiredInstanceExtensions := safeStrings(app.window.VulkanGetInstanceExtensions())
	if app.config.Debug {
		requiredInstanceExtensions = append(requiredInstanceExtensions, "VK_EXT_debug_report")
	}

	actualInstanceExtensions, err := InstanceExtensions()
	if err != nil {
		return err
	}

	log.Debug().Strs("instance_extensions", requiredInstanceExtensions).
		Msgf("loading instance extensions")

	instanceExtensions, missing := checkExisting(actualInstanceExtensions, requiredInstanceExtensions)
	if missing > 0 {
		log.Info().Msgf("missing required instance extensions during initialization")
	}

	app.extensions = instanceExtensions
	log.Info().Msgf("loaded instance extensions")

	return nil
}

func (app *Application) loadValidationLayers() error {
	var validationLayers []string
	requiredValidationLayers := safeStrings(app.validationLayers)
	actualValidationLayers, err := ValidationLayers()
	if err != nil {
		return err
	}

	log.Debug().Strs("validation_layers", validationLayers).
		Msgf("loading validation layers")

	validationLayers, missing := checkExisting(actualValidationLayers, requiredValidationLayers)
	if missing > 0 {
		log.Warn().Msgf("missing %d required validation layers during init", missing)
	}

	log.Info().Msgf("enabled instance validation layers")

	app.validationLayers = validationLayers

	return nil
}
