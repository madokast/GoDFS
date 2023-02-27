package v2

/**
单线程，内存，文件夹分级文件系统
*/

type MemFileSystem2 struct {
	fileMap map[string]*MemFileMeta2 // 文件/文件夹的全称 -> 文件源信息
}

type MemFile2 struct {
	local int64           // 当前读取位置
	fs    *MemFileSystem2 // 引用文件系统
	meta  *MemFileMeta2   // 文件元信息
}

type MemFileMeta2 struct {
	fullName    string // 文件全名
	size        int64  // 文件大小
	isDirectory bool   // 是否为文件夹

	data    []byte   // 当前文件的内容（仅当此为文件时有效）
	objects []string // 当前文件夹下的所有文件和目录，采用全名（只当此为文件夹时有效）
}

func NewMemFS() *MemFileSystem2 {
	fileMap := map[string]*MemFileMeta2{}
	fileMap["/"] = &MemFileMeta2{
		fullName:    "/",
		size:        0,
		isDirectory: true,
		data:        nil,
		objects:     []string{},
	}
	return &MemFileSystem2{fileMap: fileMap}
}
