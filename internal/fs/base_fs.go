package fs

// BaseFS 基本文件系统
type BaseFS interface {
	lock
	pathIO
	pathIOUnlock
	pathOP
	pathOPUnlock
	String() string
}

/**
注意存在有锁和无锁两种
垃圾 Go 不存在可重入锁，垃圾垃圾
*/

type pathIO interface {
	Read(path string, offset, length int64) ([]byte, error) // 读取分布式文件 path 偏移 offset 长度 length 的数据
	Write(path string, offset int64, data []byte) error     // 写入分布式文件 path 偏移 offset 数据 data
}

type pathIOUnlock interface {
	ReadUnlock(path string, offset, length int64) ([]byte, error) // 无锁版
	WriteUnlock(path string, offset int64, data []byte) error     // 无锁版
}

type pathOP interface {
	CreateFile(path string, size int64) error                         // 创建文件，指定文件大小，后期无法改变
	ListFiles(path string) (files []string, dirs []string, err error) // 列出文件夹下所有文件/路径。目录不存在不报错
	Delete(path string) error                                         // 删除文件、文件夹，如果文件夹不空则级联删除。路径不存在不会报错
	Stat(path string) (Meta, error)                                   // 获取文件元信息
	Exist(path string) (bool, error)                                  // 判断文件是否存在
}

type pathOPUnlock interface {
	CreateFileUnlock(path string, size int64) error                         // 无锁版
	ListFilesUnlock(path string) (files []string, dirs []string, err error) // 无锁版
	DeleteUnlock(path string) error                                         // 无锁版
	StatUnlock(path string) (Meta, error)                                   // 无锁版
	ExistUnlock(path string) (bool, error)                                  // 无锁版
}

type lock interface {
	RLock()
	RUnlock()
	WLock()
	WUnlock()
}
