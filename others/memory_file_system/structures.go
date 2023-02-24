package memory_file_system

/**
单线程，内存，基于一层 map 的文件系统
*/

type MemFileSystem struct {
	fileData map[string][]byte
}

type MemFile struct {
	fullName string         // 文件全名
	data     []byte         // 文件内容
	local    int            // 当前读取位置
	fs       *MemFileSystem // 引用文件系统
}

type MemFileMeta struct {
	baseName    string // 基本文件名
	fullName    string // 文件全名
	size        int64  // 文件大小
	isDirectory bool   // 是否为文件夹
}

func NewMemFS() *MemFileSystem {
	return &MemFileSystem{fileData: map[string][]byte{}}
}
