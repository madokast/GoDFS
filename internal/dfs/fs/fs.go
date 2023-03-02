package fs

import (
	"github.com/madokast/GoDFS/internal/dfs/file"
	"github.com/madokast/GoDFS/internal/dfs/node"
)

/**
分布式文件系统
*/

type Conf struct {
}

type info interface {
	AddNode(node node.Node)
	AllNodes() []node.Node
	RefreshAliveNodesAndHandCircle() // 刷新存活的 node 和 hash 环，定时调用
	HashDistribute(key string, num int) []node.Node
	String() string
}

type pathIO interface {
	Read(path string, offset, length int64) ([]byte, error) // 读取分布式文件 path 偏移 offset 长度 length 的数据
	Write(path string, offset int64, data []byte) error     // 写入分布式文件 path 偏移 offset 数据 data
}

type pathOP interface {
	CreateFile(path string, size int64) error                         // 创建文件，指定文件大小，后期无法改变
	ListFiles(path string) (files []string, dirs []string, err error) // 列出文件夹下所有文件/路径
	Delete(path string) error                                         // 删除文件、文件夹，如果文件夹不空则级联删除。路径不存在不会报错
	Stat(path string) (file.Meta, error)                              // 获取文件元信息
	Exist(path string) (bool, error)                                  // 判断文件是否存在
}

type DFS interface {
	info
	pathOP
	pathIO
}
