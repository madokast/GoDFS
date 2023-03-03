package node

import (
	"github.com/madokast/GoDFS/internal/dfs/file"
	"net/http"
)

type Info struct {
	IP      string `json:"ip,omitempty"`
	Port    uint16 `json:"port,omitempty"`
	RootDir string `json:"root_dir,omitempty"`
}

// Node 节点信息
type Node interface {
	info
	fileIO
	fileOP
	doFileOP
	sync
	ListenAndServeGo()
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Ping() bool
	DoPing(w http.ResponseWriter, r *http.Request)
}

type sync interface {
	Sync(target Node, file string) error // 从 target 同步文件到 this
	DoSync(w http.ResponseWriter, r *http.Request)
}

type info interface {
	IP() string
	Port() uint16
	RootDir() string // 文件系统根目录，OSFullPath = RootDir / FullName
	Info() *Info
	String() string
	Key() string
	Location() *file.Location
}

// nodeFileIO 发送信息到该节点处理
type fileIO interface {
	Read(path string, offset, length int64) ([]byte, error) // 节点文件读取
	Write(path string, offset int64, data []byte) error     // 节点文件写入

	DoRead(w http.ResponseWriter, r *http.Request)
	DoWrite(w http.ResponseWriter, r *http.Request)
}

// nodeFileOP 发送信息到该节点处理
type fileOP interface {
	CreateFile(path string, size int64) error                         // 创建文件，指定文件大小，后期无法改变
	ListFiles(path string) (files []string, dirs []string, err error) // 列出文件夹下所有文件/路径
	Delete(path string) error                                         // 删除文件、文件夹，如果文件夹不空则级联删除。路径不存在不会报错
	Stat(path string) (file.Meta, error)                              // 获取文件元信息
	Exist(path string) (bool, error)                                  // 判断文件是否存在
	MD5(file string) (string, error)                                  // 文件 MD5
}

type doFileOP interface {
	DoCreateFile(w http.ResponseWriter, r *http.Request)
	DoListFiles(w http.ResponseWriter, r *http.Request)
	DoDelete(w http.ResponseWriter, r *http.Request)
	DoStat(w http.ResponseWriter, r *http.Request)
	DoExist(w http.ResponseWriter, r *http.Request)
	DoMD5(w http.ResponseWriter, r *http.Request)
}
