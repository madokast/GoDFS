package memutils

import (
	"reflect"
	"unsafe"
)

// ReadBytesAs 将 src []byte 的数据指针传给到 t 中，不复制数据
// 也许存在内存泄漏，合理的使用方法如下
//
//	    var p *Pointer
//		ReadBytesAs(bytes, &p)
//		return *p
func ReadBytesAs[T interface{}](src []byte, t **T) {
	srcHelper := (*reflect.SliceHeader)(unsafe.Pointer(&src))
	*t = (*T)(unsafe.Pointer(srcHelper.Data))
}

// WriteBytesTo 将 t 的内存复制到 dst 中。复制大小为 dst 的长度
func WriteBytesTo[T interface{}](dst []byte, t *T) {
	srcHeader := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(t)),
		Cap:  len(dst),
		Len:  len(dst),
	}

	src := *(*[]byte)(unsafe.Pointer(&srcHeader))
	copy(dst, src)
}

// WriteAsBytes 将指针 t 指向的 size 区域的内存复制到 []byte 并返回
func WriteAsBytes[T interface{}](t *T, size uint32) []byte {
	b := make([]byte, size)
	src := *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(t)),
		Cap:  int(size),
		Len:  int(size),
	}))
	copy(b, src)
	return b
}
