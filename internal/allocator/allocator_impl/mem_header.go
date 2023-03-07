package allocator_impl

import (
	"fmt"
	"github.com/madokast/GoDFS/internal/allocator/memutils"
	"github.com/madokast/GoDFS/utils/logger"
	"strconv"
	"unsafe"
)

// MemHeader 内存头信息，32 位
// 第一位表示内存是否释放。1 正在使用，0 释放
// 后 31 位表示内存大小
type MemHeader struct {
	freeFlagSize uint32
}

var MemHeaderSz = uint32(unsafe.Sizeof(MemHeader{}))

func CreateHeader(size uint32) MemHeader {
	if size&0x8000_0000 == 0x8000_0000 {
		panic("Size too large " + strconv.Itoa(int(size)))
	}
	return MemHeader{freeFlagSize: size | 0x8000_0000}
}

func (h *MemHeader) Size() uint32 {
	return h.freeFlagSize & 0x7FFF_FFFF
}

func (h *MemHeader) IsFree() bool {
	return h.freeFlagSize&0x8000_0000 == 0x0
}

func (h *MemHeader) WriteHeaderTo(data []byte) {
	memutils.WriteBytesTo(data[:MemHeaderSz], h)
}

func (h *MemHeader) Free() {
	if h.IsFree() {
		logger.Error("Double free!!")
	}
	h.freeFlagSize &= 0x7FFF_FFFF
}

func (h *MemHeader) String() string {
	use := "USE"
	if h.IsFree() {
		use = "FREE"
	}
	return fmt.Sprintf("MH(%d)%s", h.Size(), use)
}

func ReadHeader(data []byte) MemHeader {
	var h *MemHeader
	memutils.ReadBytesAs(data[:MemHeaderSz], &h)
	return *h
}
