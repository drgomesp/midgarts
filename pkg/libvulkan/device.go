package libvulkan

import (
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/veandco/go-sdl2/sdl"
	vk "github.com/vulkan-go/vulkan"
)

type Device struct {
	config Config

	instance vk.Instance
	gpu      vk.PhysicalDevice
	device   vk.Device
	surface  vk.Surface

	extensions, validationLayers []string

	gpuProperties    vk.PhysicalDeviceProperties
	memoryProperties vk.PhysicalDeviceMemoryProperties

	graphicsQueueIndex uint32
	presentQueueIndex  uint32
	presentQueue       vk.Queue
	graphicsQueue      vk.Queue
}

func NewDevice(config Config, extensions []string, validationLayers []string, window *sdl.Window) (*Device, error) {
	d := &Device{
		config:           config,
		extensions:       extensions,
		validationLayers: validationLayers,
	}

	if err := d.createInstance(extensions, validationLayers); err != nil {
		return nil, err
	}

	if err := d.loadPhysicalDevice(); err != nil {
		return nil, err
	}

	if err := d.loadLogicalDevice(); err != nil {
		return nil, err
	}

	if err := d.loadExtensions(); err != nil {
		return nil, err
	}

	if err := d.loadSurface(window); err != nil {
		return nil, err
	}

	// Get queue family properties
	var queueCount uint32
	vk.GetPhysicalDeviceQueueFamilyProperties(d.gpu, &queueCount, nil)
	queueProperties := make([]vk.QueueFamilyProperties, queueCount)
	vk.GetPhysicalDeviceQueueFamilyProperties(d.gpu, &queueCount, queueProperties)
	if queueCount == 0 { // probably should try another GPU
		return nil, errors.New("vulkan error: no queue families found on GPU 0")
	}

	return d, nil
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
			PEngineName:        safeString("libVulkan"),
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

func (d *Device) loadPhysicalDevice() error {
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

	log.Debug().Msgf("selected GPU device [%s]", d.gpuProperties.DeviceName)

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

func (d *Device) loadSurface(window *sdl.Window) error {
	surfPtr, err := window.VulkanCreateSurface(d.instance)
	if err != nil {
		return err
	}

	d.surface = vk.SurfaceFromPointer(uintptr(surfPtr))
	return nil
}

func (d *Device) loadLogicalDevice() error {
	// Get queue family properties
	var queueCount uint32
	vk.GetPhysicalDeviceQueueFamilyProperties(d.gpu, &queueCount, nil)
	queueProperties := make([]vk.QueueFamilyProperties, queueCount)
	vk.GetPhysicalDeviceQueueFamilyProperties(d.gpu, &queueCount, queueProperties)
	if queueCount == 0 { // probably should try another GPU
		return errors.New("vulkan error: no queue families found on GPU 0")
	}

	// Find a suitable queue family for the target Vulkan mode
	var graphicsFound bool
	var presentFound bool
	var separateQueue bool
	for i := uint32(0); i < queueCount; i++ {
		var (
			required        vk.QueueFlags
			supportsPresent vk.Bool32
			needsPresent    bool
		)
		if graphicsFound {
			// looking for separate present queue
			separateQueue = true
			vk.GetPhysicalDeviceSurfaceSupport(d.gpu, i, d.surface, &supportsPresent)
			if supportsPresent.B() {
				d.presentQueueIndex = i
				presentFound = true
				break
			}
		}

		required |= vk.QueueFlags(vk.QueueGraphicsBit)
		queueProperties[i].Deref()

		if queueProperties[i].QueueFlags&required != 0 {
			if !needsPresent || (needsPresent && supportsPresent.B()) {
				d.graphicsQueueIndex = i
				graphicsFound = true
				break
			} else if needsPresent {
				d.graphicsQueueIndex = i
				graphicsFound = true
			}
		}
	}

	if separateQueue && !presentFound {
		return errors.New("vulkan error: could not found separate queue with present capabilities")
	}
	if !graphicsFound {
		return errors.New("vulkan error: could not find a suitable queue family for the target Vulkan mode")
	}

	// Create a Vulkan device
	queueInfos := []vk.DeviceQueueCreateInfo{{
		SType:            vk.StructureTypeDeviceQueueCreateInfo,
		QueueFamilyIndex: d.graphicsQueueIndex,
		QueueCount:       1,
		PQueuePriorities: []float32{1.0},
	}}

	if separateQueue {
		queueInfos = append(queueInfos, vk.DeviceQueueCreateInfo{
			SType:            vk.StructureTypeDeviceQueueCreateInfo,
			QueueFamilyIndex: d.presentQueueIndex,
			QueueCount:       1,
			PQueuePriorities: []float32{1.0},
		})
	}

	var device vk.Device
	ret := vk.CreateDevice(d.gpu, &vk.DeviceCreateInfo{
		SType:                   vk.StructureTypeDeviceCreateInfo,
		QueueCreateInfoCount:    uint32(len(queueInfos)),
		PQueueCreateInfos:       queueInfos,
		EnabledExtensionCount:   uint32(len(d.extensions)),
		PpEnabledExtensionNames: d.extensions,
		EnabledLayerCount:       uint32(len(d.validationLayers)),
		PpEnabledLayerNames:     d.validationLayers,
	}, nil, &device)
	if err := VkError(ret); err != nil {
		return err
	}

	d.device = device

	//d.context.device = device
	//app.VulkanInit(p.context)

	var queue vk.Queue
	vk.GetDeviceQueue(d.device, d.graphicsQueueIndex, 0, &queue)
	d.graphicsQueue = queue

	if d.config.Mode.Has(VulkanPresent) {
		if separateQueue {
			var presentQueue vk.Queue
			vk.GetDeviceQueue(d.device, d.presentQueueIndex, 0, &presentQueue)
			d.presentQueue = presentQueue
		}
		//p.context.preparePresent()

		//dimensions := &SwapchainDimensions{
		//	// some default preferences here
		//	Width: 640, Height: 480,
		//	Format: vk.FormatB8g8r8a8Unorm,
		//}
		//if iface, ok := app.(ApplicationSwapchainDimensions); ok {
		//	dimensions = iface.VulkanSwapchainDimensions()
		//}
		//p.context.prepareSwapchain(p.gpu, p.surface, dimensions)
	}
	//if iface, ok := app.(ApplicationContextPrepare); ok {
	//	p.context.SetOnPrepare(iface.VulkanContextPrepare)
	//}
	//if iface, ok := app.(ApplicationContextCleanup); ok {
	//	p.context.SetOnCleanup(iface.VulkanContextCleanup)
	//}
	//if iface, ok := app.(ApplicationContextInvalidate); ok {
	//	p.context.SetOnInvalidate(iface.VulkanContextInvalidate)
	//}
	//if mode.Has(VulkanPresent) {
	//	p.context.prepare(false)
	//}

	return nil
}