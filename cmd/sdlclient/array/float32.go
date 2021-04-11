package array

import "unsafe"

type Float32 []float32

func NewArrayFloat32(size, cap int) Float32 {
	return make([]float32, size, cap)
}

func (f *Float32) Append(v ...float32) {
	*f = append(*f, v...)
}

func (f *Float32) Bytes() int {
	return len(*f) * int(unsafe.Sizeof(float32(0)))
}

func (f Float32) Float32Array() []float32 {
	return (*[1 << 27]float32)(unsafe.Pointer(&f[0]))[:len(f)]
}
