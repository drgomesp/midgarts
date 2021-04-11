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
	return (*[1 << 27]uint32)(unsafe.Pointer(&u[0]))[:len(u)]
}

func (u *Uint32) Size() int {
	return len(*u)
}
