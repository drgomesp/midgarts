package libvulkan

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	vk "github.com/vulkan-go/vulkan"
)

type Device struct {
	config Config

	instance vk.Instance
	gpu      vk.PhysicalDevice
	device   vk.Device

	gpuProperties    vk.PhysicalDeviceProperties
	memoryProperties vk.PhysicalDeviceMemoryProperties
}

func NewDevice(config Config, extensions []string, validationLayers []string) (*Device, error) {
	device := &Device{
		config: config,
	}

	if err := device.createInstance(extensions, validationLayers); err != nil {
		return nil, err
	}

	if err := device.selectPhysicalDevice(); err != nil {
		return nil, err
	}

	if err := device.loadExtensions(); err != nil {
		return nil, err
	}

	return device, nil
}

func (d *Device) createInstance(extensions []string, validationLayers []string) error {
	var instance vk.Instance
	ret := vk.CreateInstance(&vk.InstanceCreateInfo{
		SType: vk.StructureTypeInstanceCreateInfo,
		PApplicationInfo: &vk.ApplicationInfo{
			SType:              vk.StructureTypeApplicationInfo,
			ApiVersion:         uint32(d.config.APIVersion),
			ApplicationVersion: uint32(d.config.AppVersion),
			PApplicationName:   safeString(d.config.AppName),
			PEngineName:        "kurwamuc\x00",
		},
		EnabledExtensionCount:   uint32(len(extensions)),
		PpEnabledExtensionNames: extensions,
		EnabledLayerCount:       uint32(len(validationLayers)),
		PpEnabledLayerNames:     validationLayers,
	}, nil, &instance)
	if err := VkError(ret); err != nil {
		return err
	}

	d.instance = instance
	err := vk.InitInstance(instance)

	if err != nil {
		return err
	}

	return nil
}

func (d *Device) selectPhysicalDevice() error {
	var gpuCount uint32
	ret := vk.EnumeratePhysicalDevices(d.instance, &gpuCount, nil)
	if err := VkError(ret); err != nil {
		return err
	}

	if gpuCount == 0 {
		return errors.New("no GPU devices found")
	}

	gpus := make([]vk.PhysicalDevice, gpuCount)
	ret = vk.EnumeratePhysicalDevices(d.instance, &gpuCount, gpus)
	if err := VkError(ret); err != nil {
		return err
	}

	d.gpu = gpus[0]
	vk.GetPhysicalDeviceProperties(d.gpu, &d.gpuProperties)
	d.gpuProperties.Deref()
	vk.GetPhysicalDeviceMemoryProperties(d.gpu, &d.memoryProperties)
	d.memoryProperties.Deref()

	log.Debug().Msgf("selected GPU device %s]", d.gpuProperties.DeviceName)

	return nil
}

func (d *Device) loadExtensions() error {
	requiredDeviceExtensions := safeStrings(d.config.DeviceExtensions)

	var count uint32
	ret := vk.EnumerateDeviceExtensionProperties(d.gpu, "", &count, nil)
	if err := VkError(ret); err != nil {
		return err
	}

	var actualDeviceExtensions []string
	list := make([]vk.ExtensionProperties, count)
	ret = vk.EnumerateDeviceExtensionProperties(d.gpu, "", &count, list)
	if err := VkError(ret); err != nil {
		return err
	}

	for _, ext := range list {
		ext.Deref()
		actualDeviceExtensions = append(actualDeviceExtensions, vk.ToString(ext.ExtensionName[:]))
	}

	deviceExtensions, missing := checkExisting(actualDeviceExtensions, requiredDeviceExtensions)
	if missing > 0 {
		log.Warn().Msgf("missing %d required device extensions during init", missing)
	}

	log.Debug().Strs("device_extensions", deviceExtensions).Msgf("device extensions")

	log.Info().Msg("enabled device extensions")

	return nil
}
