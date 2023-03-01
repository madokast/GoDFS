package dfs

import "net/http"

type NodeConf struct {
	IP      string
	Port    uint16
	RootDir string
}

// Node 节点信息
type Node interface {
	nodeInfo
	nodeFileIO
	nodeDoFileIO
	nodeFileOP
	nodeDoFileOP
	ListenAndServeGo()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Ping() bool
	DoPing(w http.ResponseWriter, r *http.Request)
}

type nodeInfo interface {
	IP() string
	Port() uint16
	RootDir() string // 文件系统根目录，OSFullPath = RootDir / FullName
	String() string
}

// nodeFileIO 发送信息到该节点处理
type nodeFileIO interface {
	Lock(path string) error                                 // 节点文件加锁
	Unlock(path string) error                               // 节点文件解锁
	SetVersion(path string) error                           // 设定版本
	Version(path string) (int64, error)                     // 节点文件版本。需要有 header 存储锁信息和版本
	Read(path string, offset, length int64) ([]byte, error) // 节点文件读取
	Write(path string, offset int64, data []byte) error     // 节点文件写入
}

// nodeDoFileIO 节点处理逻辑，在
type nodeDoFileIO interface {
	DoLock(w http.ResponseWriter, r *http.Request)
	DoUnlock(w http.ResponseWriter, r *http.Request)
	DoSetVersion(w http.ResponseWriter, r *http.Request)
	DoVersion(w http.ResponseWriter, r *http.Request)
	DoRead(w http.ResponseWriter, r *http.Request)
	DoWrite(w http.ResponseWriter, r *http.Request)
}

// nodeFileOP 发送信息到该节点处理
type nodeFileOP interface {
	CreateFile(path string, size int64) error                         // 创建文件，指定文件大小，后期无法改变
	MkdirAll(path string)                                             // 创建文件夹，支持级联
	ListFiles(path string) (files []string, dirs []string, err error) // 列出文件夹下所有文件/路径
	Delete(path string) error                                         // 删除文件、文件夹，如果文件夹不空则级联删除。路径不存在不会报错
	Stat(path string) (FileMeta, error)                               // 获取文件元信息
	Exist(path string) (bool, error)                                  // 判断文件是否存在
}

type nodeDoFileOP interface {
	DoCreateFile(w http.ResponseWriter, r *http.Request)
	DoMkdirAll(w http.ResponseWriter, r *http.Request)
	DoListFiles(w http.ResponseWriter, r *http.Request)
	DoDelete(w http.ResponseWriter, r *http.Request)
	DoStat(w http.ResponseWriter, r *http.Request)
	DoExist(w http.ResponseWriter, r *http.Request)
}
