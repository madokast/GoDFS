package nodeimpl

import (
	"fmt"
	"github.com/madokast/GoDFS/internal/dfs/file"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"net/http"
)

const (
	baseApi       = "/node"
	pingApi       = baseApi + "/alive"
	readFileApi   = baseApi + "/file/read"
	writeFileApi  = baseApi + "/file/write"
	createFileApi = baseApi + "/file/create"
	md5Api        = baseApi + "/file/md5"
	syncApi       = baseApi + "/file/sync"
	deletePathApi = baseApi + "/path/delete"
	statPathApi   = baseApi + "/path/stat"
	existPathApi  = baseApi + "/path/exist"
	listFilesApi  = baseApi + "/listfiles"
)

type Impl struct {
	ip      string
	port    uint16
	rootDir string
	router  map[string]web.HandleFunc
}

func New(conf *node.Info) node.Node {
	n := &Impl{ip: conf.IP, port: conf.Port, rootDir: conf.RootDir}
	n.registerRouter()
	return n
}

func (n *Impl) registerRouter() {
	n.router = map[string]web.HandleFunc{}
	n.router[pingApi] = n.DoPing
	// node_file-io.go
	n.router[readFileApi] = n.DoRead
	n.router[writeFileApi] = n.DoWrite
	// node_file-op.go
	n.router[createFileApi] = n.DoCreateFile
	n.router[deletePathApi] = n.DoDelete
	n.router[md5Api] = n.DoMD5
	n.router[syncApi] = n.DoSync
	// node_path-check.go
	n.router[statPathApi] = n.DoStat
	n.router[existPathApi] = n.DoExist
	// node_dir-op.go
	n.router[listFilesApi] = n.DoListFiles
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

func (n *Impl) Location() *file.Location {
	return &file.Location{
		IP:      n.ip,
		Port:    n.port,
		RootDir: n.rootDir,
	}
}

func (n *Impl) String() string {
	return fmt.Sprintf("Node %s:%d root %s", n.ip, n.port, n.rootDir)
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

func (n *Impl) Ping() bool {
	ret := &web.Response[*web.NullResponse]{}
	err := httputils.PostJson(n.ip, n.port, pingApi, &web.NullRequest{}, ret)
	if err != nil {
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
	handle, ok := n.router[r.URL.Path]
	if ok {
		handle(w, r)
	} else {
		panic("Unknown url " + r.URL.Path)
	}
}

func (n *Impl) ListenAndServeGo() {
	logger.Info("Start node", n.Location())
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", n.ip, n.port), n)
		if err != nil {
			logger.Error(err)
		}
	}()
}
