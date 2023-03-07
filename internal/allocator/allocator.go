package allocator

import (
	"fmt"
	"github.com/madokast/GoDFS/internal/allocator/memutils"
	"github.com/madokast/GoDFS/internal/fs/write_callback"
	"unsafe"
)

/**
allocator 基于底层的文件提供内存分配器
每个内存分配器，基于一个本地文件夹，例如 /root，文件夹内创建固定大小的若干文件对外分配内存
*/

// Pointer 指针
type Pointer struct {
	BlockId     uint32 // 可以看作文件编号
	BlockOffset uint32 // 可以看作文件内偏移。blockId 和 blockOffset 一起完成寻址
}

type CacheData struct {
	Data  []byte
	WcObj *write_callback.Entry
}

type Allocator interface {
	Allocate(size uint32, newBlock bool) (Pointer, error) // 分配内存
	Free(p Pointer) error                                 // free 内存
	Read(p Pointer, reader func([]byte)) error            // 读取
	ReadString(p Pointer) (string, error)
	ReadBytes(p Pointer) ([]byte, error)
	Write(p Pointer, offset uint32, data []byte) error // 写入
	String() string
}

var PointerSz = uint32(unsafe.Sizeof(Pointer{}))

// ReadPointer 从 ptrBytes 中读取 Pointer
func ReadPointer(ptrBytes []byte) Pointer {
	var p *Pointer
	memutils.ReadBytesAs(ptrBytes, &p)
	return *p
}

func (p Pointer) ToBytes() []byte {
	return memutils.WriteAsBytes(&p, PointerSz)
}

func (p Pointer) String() string {
	return fmt.Sprintf("%d:%d", p.BlockId, p.BlockOffset)
}
