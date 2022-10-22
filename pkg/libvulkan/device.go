package libvulkan

import vk "github.com/vulkan-go/vulkan"

type Device struct {
	config Config

	instance vk.Instance
	gpu      vk.PhysicalDevice
	device   vk.Device
}

func NewDevice(config Config, extensions []string, validationLayers []string) (*Device, error) {
	device := &Device{
		config: config,
	}

	if err := device.createInstance(extensions, validationLayers); err != nil {
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
