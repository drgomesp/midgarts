package libvulkan

import vk "github.com/vulkan-go/vulkan"

var (
	DefaultVulkanAppVersion = vk.MakeVersion(1, 0, 0)
	DefaultVulkanAPIVersion = vk.MakeVersion(1, 0, 0)
)

type Config struct {
	AppName          string
	APIVersion       vk.Version
	AppVersion       vk.Version
	ValidationLayers []string
	DeviceExtensions []string
}

func DefaultConfig() Config {
	return Config{
		AppName:          "libvulkan-app",
		APIVersion:       vk.Version(DefaultVulkanAPIVersion),
		AppVersion:       vk.Version(DefaultVulkanAppVersion),
		ValidationLayers: []string{
			// "VK_LAYER_GOOGLE_threading",
			// "VK_LAYER_LUNARG_parameter_validation",
			// "VK_LAYER_LUNARG_object_tracker",
			// "VK_LAYER_LUNARG_core_validation",
			// "VK_LAYER_LUNARG_api_dump",
			// "VK_LAYER_LUNARG_swapchain",
			// "VK_LAYER_GOOGLE_unique_objects",
		},
		DeviceExtensions: []string{
			"VK_KHR_swapchain",
		},
	}
}
