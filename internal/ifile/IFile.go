package ifile

import (
	"io"
)

// FileIO 文件读写操作
type FileIO interface {
	io.Closer                             // 关闭文件IO
	io.Reader                             // 读文件
	io.Writer                             // 写文件
	Offset(offset int64) error            // 移动读写位置
	ReadString(limit int) (string, error) // 读字符串
}

// FileMeta 文件元信息
type FileMeta interface {
	BaseName() string  // 基本文件名
	FullName() string  // 文件全名
	Size() int64       // 文件大小
	IsDirectory() bool // 是否为文件夹
}

// FileManager 文件管理
type FileManager interface {
	CreateFile(name string) (FileIO, error) // 创建文件
	OpenFile(name string) (FileIO, error)   // 打开文件用于读写
	DeleteFile(name string) error           // 删除文件
	RenameFile(name, newName string) error  // 改名文件
}

// DirectoryManager 文件夹管理
type DirectoryManager interface {
	MakeDirectory(name string) error              // 创建文件夹
	MakeDirectories(name string) error            // 创建文件夹，如有需要创建父文件夹
	DeleteDirectory(name string) error            // 删除空文件夹
	DeleteDirectories(name string) error          // 删除文件夹，如有必要删除文件夹内所有内容
	ReadDirectory(dir string) ([]FileMeta, error) // 列出文件夹下文件
	RenameDirectory(name, newName string) error   // 改名文件夹
}

// FileSystem 文件系统
type FileSystem interface {
	FileManager
	DirectoryManager
	Stat(name string) (FileMeta, error) // 获取元信息
}
