package v2

import (
	"strconv"
	"sync"
	"sync/atomic"
)

/**
单线程，内存，文件夹分级文件系统
*/

// 文件模式。处于读模式下时，前 16 位记录打开次数
const (
	closeMode int32 = 0
	readMode  int32 = 1
	writeMode int32 = 3
)

type MemFileSystem2 struct {
	fileMap map[string]*MemFileMeta2 // 文件/文件夹的全称 -> 文件源信息
	fsLock  sync.Mutex               // 文件系统的全局锁，进行 stat、创建/删除等时候使用
}

type MemFileDescription2 struct {
	local  int64         // 当前读取位置
	closed bool          // 文件是否已经关闭
	meta   *MemFileMeta2 // 文件元信息
	fdLock sync.Mutex    // 因为有一个读写位置的概念，而且不能关闭两次，所以都要加锁
}

type MemFileMeta2 struct {
	fullName    string // 文件全名
	size        int64  // 文件大小
	isDirectory bool   // 是否为文件夹

	data         []byte              // 当前文件的内容（当这个 meta 是文件时才使用）
	objects      map[string]struct{} // 当前文件夹下的所有文件和目录，采用全名（当这个 meta 是文件夹时才使用）
	openMode     atomic.Int32        // 文件的读写模式。只有当这个 meta 是文件时才使用
	openModeLock sync.RWMutex        // 读写 openMode 的锁
}

func NewMemFS() *MemFileSystem2 {
	fileMap := map[string]*MemFileMeta2{}
	fileMap["/"] = &MemFileMeta2{
		fullName:    "/",
		size:        0,
		isDirectory: true,
		data:        nil,
		objects:     map[string]struct{}{},
	}
	return &MemFileSystem2{fileMap: fileMap}
}

func modeString(openMode int32) string {
	if openMode == closeMode {
		return "closed"
	}
	if openMode == writeMode {
		return "write"
	}
	if isReadMode(openMode) {
		readRef := openMode >> 16
		return "read[" + strconv.Itoa(int(readRef)) + "]"
	}
	panic(openMode)
}

func isReadMode(openMode int32) bool {
	openMode &= 0x0000_FFFF
	return openMode == readMode
}

func canReadMode(openMode int32) bool {
	return (openMode & readMode) == readMode
}
