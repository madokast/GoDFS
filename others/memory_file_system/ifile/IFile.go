package ifile

import (
	"io"
)

// FileIO 文件读写操作。文件大小固定，不能更改
type FileIO interface {
	io.Closer                             // 关闭文件IO
	io.Reader                             // 读文件
	io.Writer                             // 写文件
	Offset(offset int64) error            // 移动读写位置
	ReadString(limit int) (string, error) // 读字符串
}

// FileMeta 文件元信息
type FileMeta interface {
	FullName() string  // 文件全名
	Size() int64       // 文件大小
	IsDirectory() bool // 是否为文件夹
	String() string
}

// FileManager 文件管理，所有名称必须是全称 /a/b/c 形式
type FileManager interface {
	CreateFile(name string, size int64) error         // 创建文件
	OpenFile(name string, write bool) (FileIO, error) // 打开文件。write 标记只读还是读写
	DeleteFile(name string) error                     // 删除文件
}

// DirectoryManager 文件夹管理，所有名称必须是全称 /a/b/c 形式
type DirectoryManager interface {
	MakeDirectory(name string) error                                     // 创建文件夹
	MakeDirectories(name string) error                                   // 创建文件夹，如有需要创建父文件夹
	DeleteDirectory(name string) error                                   // 删除空文件夹
	DeleteDirectoryAll(name string) error                                // 删除文件夹，如有必要删除文件夹内所有内容
	ReadDirectory(dir string) (files []string, dirs []string, err error) // 列出文件夹下文件
}

// FileSystem 文件系统
type FileSystem interface {
	FileManager
	DirectoryManager
	Stat(name string) (FileMeta, error) // 获取元信息
	Exist(name string) bool             // 文件是否存在
}
