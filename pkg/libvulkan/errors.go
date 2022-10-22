package libvulkan

import (
	"fmt"
	"runtime"

	vk "github.com/vulkan-go/vulkan"
)

func isError(ret vk.Result) bool {
	return ret != vk.Success
}

func VkError(ret vk.Result) error {
	if ret != vk.Success {
		pc, _, _, ok := runtime.Caller(0)
		if !ok {
			return fmt.Errorf("vulkan error: %s (%d)",
				vk.Error(ret).Error(), ret)
		}
		frame := newStackFrame(pc)
		return fmt.Errorf("vulkan error: %s (%d) on %s",
			vk.Error(ret).Error(), ret, frame.String())
	}
	return nil
}

func orPanic(err error, finalizers ...func()) {
	if err != nil {
		for _, fn := range finalizers {
			fn()
		}
		panic(err)
	}
}

func checkErr(err *error) {
	if v := recover(); v != nil {
		*err = fmt.Errorf("%+v", v)
	}
}
