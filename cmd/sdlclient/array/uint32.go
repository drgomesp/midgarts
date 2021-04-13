package array

import "unsafe"

type Uint32 []uint32

func NewArrayUint32(size, cap int) Uint32 {
	return make([]uint32, size, cap)
}

func (u *Uint32) Append(v ...uint32) {
	*u = append(*u, v...)
}

func (u *Uint32) Bytes() int {
	return len(*u) * int(unsafe.Sizeof(uint32(0)))
}

func (u Uint32) ToUint32() []uint32 {
	return u[:]
}

func (u Uint32) Size() int32 {
	return int32(len(u))
}
