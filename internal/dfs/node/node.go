package node

import (
	"github.com/madokast/GoDFS/internal/dfs/dfile"
	"github.com/madokast/GoDFS/internal/fs"
	"github.com/madokast/GoDFS/internal/fs/write_callback"
	"net/http"
)

const (
	baseApi          = "/node"
	PingApi          = baseApi + "/alive"
	ReadFileApi      = baseApi + "/file/read"
	WriteFileApi     = baseApi + "/file/write"
	CreateFileApi    = baseApi + "/file/create"
	Md5Api           = baseApi + "/file/md5"
	SyncApi          = baseApi + "/file/sync"
	WriteCallBackApi = baseApi + "/file/write-callback"
	DeletePathApi    = baseApi + "/path/delete"
	StatPathApi      = baseApi + "/path/stat"
	ExistPathApi     = baseApi + "/path/exist"
	ListFilesApi     = baseApi + "/list-files"
)

type Info struct {
	IP      string `json:"ip,omitempty"`
	Port    uint16 `json:"port,omitempty"`
	RootDir string `json:"root_dir,omitempty"`
}

// Node 节点信息
type Node interface {
	fs.BaseFS // 发送信息到该节点处理
	info
	doFileIO
	fileOP
	doFileOP
	sync
	write_callback.WriteCallBack
	ListenAndServeGo()                                // 启动 node 的 rpc 服务
	ServeHTTP(w http.ResponseWriter, r *http.Request) // 实现 http.Handler 接口
	Close()                                           // node rpc 服务下线，一般只用于测试
}

type sync interface {
	Ping() bool // 检查 node 是否存活
	DoPing(w http.ResponseWriter, r *http.Request)

	Sync(src Node, file string) error // 从 src 同步文件到 this
	DoSync(w http.ResponseWriter, r *http.Request)
}

type info interface {
	IP() string
	Port() uint16
	RootDir() string // 文件系统根目录，OSFullPath = RootDir / FullName
	Info() *Info
	String() string
	Key() string
	Location() *dfile.Location
	IsLocalService() bool // 是否为本地 rpc，完成一些本地回调
}

type doFileIO interface {
	DoRead(w http.ResponseWriter, r *http.Request)
	DoWrite(w http.ResponseWriter, r *http.Request)
}

// nodeFileOP 发送信息到该节点处理
type fileOP interface {
	MD5(file string) (string, error)                    // 文件 MD5
	ForAllFile(path string, consumer func(file string)) // 目录下文件深度遍历。如果传入的是文件，则直接传给 consumer
}

type doFileOP interface {
	DoCreateFile(w http.ResponseWriter, r *http.Request)
	DoListFiles(w http.ResponseWriter, r *http.Request)
	DoDelete(w http.ResponseWriter, r *http.Request)
	DoStat(w http.ResponseWriter, r *http.Request)
	DoExist(w http.ResponseWriter, r *http.Request)
	DoMD5(w http.ResponseWriter, r *http.Request)
}
