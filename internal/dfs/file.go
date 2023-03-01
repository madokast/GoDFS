package dfs

import "unsafe"

type FileHeader struct {
	Version  int64
	LockFlag int32
}

var FileHeaderSz = int64(unsafe.Sizeof(FileHeader{}))

// FileMeta 文件信息
type FileMeta interface {
	FullName() string  // 文件全名
	Size() int64       // 文件大小
	IsDirectory() bool // 是否为文件夹
	Locations() []Node // 文件位置
	Version() int64
	String() string
}
