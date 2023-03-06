package nodeimpl

import (
	"fmt"
	"github.com/madokast/GoDFS/internal/dfs/dfile"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/dfs/node/write_callback"
	"github.com/madokast/GoDFS/internal/fs"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"net/http"
	"strings"
)

type Impl struct {
	ip               string
	port             uint16
	rootDir          string
	router           map[string]web.HandleFunc
	closed           bool // rpc 服务是否关闭。一般用于测试
	localService     bool // 是否为本地 rpc，完成一些本地回调
	fs.WriteCallBack      // 写操作监听回调
}

func New(conf *node.Info) node.Node {
	n := &Impl{
		ip:           conf.IP,
		port:         conf.Port,
		rootDir:      conf.RootDir,
		closed:       false,
		localService: false,
	}
	n.WriteCallBack = write_callback.New(n)
	n.registerRouter()
	return n
}

func (n *Impl) registerRouter() {
	n.router = map[string]web.HandleFunc{}
	n.router[node.PingApi] = n.DoPing
	// node_file-io.go
	n.router[node.ReadFileApi] = n.DoRead
	n.router[node.WriteFileApi] = n.DoWrite
	// node_file-op.go
	n.router[node.CreateFileApi] = n.DoCreateFile
	n.router[node.DeletePathApi] = n.DoDelete
	n.router[node.Md5Api] = n.DoMD5
	n.router[node.SyncApi] = n.DoSync
	// node_path-check.go
	n.router[node.StatPathApi] = n.DoStat
	n.router[node.ExistPathApi] = n.DoExist
	// node_dir-op.go
	n.router[node.ListFilesApi] = n.DoListFiles
	// writeCallBackApi
	n.router[node.WriteCallBackApi] = n.DoWriteCallback
}

func (n *Impl) IP() string {
	return n.ip
}

func (n *Impl) Port() uint16 {
	return n.port
}

func (n *Impl) RootDir() string {
	return n.rootDir
}

func (n *Impl) Location() *dfile.Location {
	return &dfile.Location{
		IP:      n.ip,
		Port:    n.port,
		RootDir: n.rootDir,
	}
}

func (n *Impl) String() string {
	return n.Key()
}

func (n *Impl) Info() *node.Info {
	return &node.Info{
		IP:      n.ip,
		Port:    n.port,
		RootDir: n.rootDir,
	}
}

func (n *Impl) Key() string {
	return fmt.Sprintf("%s:%d", n.ip, n.port)
}

func (n *Impl) baseUrl() string {
	return fmt.Sprintf("http://%s:%d", n.ip, n.port)
}

func (n *Impl) IsLocalService() bool {
	return n.localService
}

func (n *Impl) Ping() bool {
	ret := &web.Response[*web.NullResponse]{}
	err := httputils.PostJson(n.ip, n.port, node.PingApi, &web.NullRequest{}, ret)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		logger.Error(err)
		return false
	}
	return ret.Msg == web.SuccessMsg
}

func (n *Impl) DoPing(w http.ResponseWriter, r *http.Request) {
	httputils.HandleJson(w, r, &web.BaseResponse{}, func(req *web.BaseResponse) (*web.NullResponse, error) {
		return &web.NullResponse{}, nil
	})
}

func (n *Impl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if n.closed {
		w.WriteHeader(500)
		return
	}

	handle, ok := n.router[r.URL.Path]
	if ok {
		handle(w, r)
	} else {
		panic("Unknown url " + r.URL.Path)
	}
}

func (n *Impl) ListenAndServeGo() {
	logger.Info("Start node", n.Location())
	n.closed = false
	n.localService = true
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", n.ip, n.port), n)
		if err != nil {
			logger.Error(err)
		}
	}()
}

func (n *Impl) Close() {
	n.closed = true
}
