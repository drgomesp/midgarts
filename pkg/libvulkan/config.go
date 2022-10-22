package libvulkan

import vk "github.com/vulkan-go/vulkan"

var (
	DefaultVulkanAppVersion = vk.MakeVersion(1, 0, 0)
	DefaultVulkanAPIVersion = vk.MakeVersion(1, 0, 0)
)

type Config struct {
	APIVersion vk.Version
	AppVersion vk.Version
}

func DefaultConfig() *Config {
	return &Config{
		APIVersion: vk.Version(DefaultVulkanAPIVersion),
		AppVersion: vk.Version(DefaultVulkanAppVersion),
	}
}
