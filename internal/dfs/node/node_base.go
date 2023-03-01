package node

import (
	"fmt"
	"github.com/madokast/GoDFS/internal/dfs"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"net/http"
)

const (
	baseApi = "/node"
	pingApi = baseApi + "/alive"
	readApi = baseApi + "/read"
)

type node struct {
	ip      string
	port    uint16
	rootDir string
}

func New(conf *dfs.NodeConf) dfs.Node {
	return &node{ip: conf.IP, port: conf.Port, rootDir: conf.RootDir}
}

func (n *node) IP() string {
	return n.ip
}

func (n *node) Port() uint16 {
	return n.port
}

func (n *node) RootDir() string {
	return n.rootDir
}

func (n *node) String() string {
	return fmt.Sprintf("node %s:%d root-dir %s", n.ip, n.port, n.rootDir)
}

func (n *node) baseUrl() string {
	return fmt.Sprintf("http://%s:%d", n.ip, n.port)
}

func (n *node) Ping() bool {
	ret := &web.Response[*web.NullResponse]{}
	err := httputils.PostJson(n.ip, n.port, pingApi, &web.NullRequest{}, ret)
	if err != nil {
		logger.Error(err)
		return false
	}
	return ret.Msg == web.SuccessMsg
}

func (n *node) DoPing(w http.ResponseWriter, r *http.Request) {
	httputils.HandleJson(w, r, &web.BaseResponse{}, func(req *web.BaseResponse) (*web.NullResponse, error) {
		return &web.NullResponse{}, nil
	})
}

func (n *node) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case pingApi:
		n.DoPing(w, r)
	case readApi:
		n.DoRead(w, r)
	default:
		logger.Error("Unknown url", r.URL.Path)
	}
}

func (n *node) ListenAndServeGo() {
	go func() {
		err := http.ListenAndServe(fmt.Sprintf("%s:%d", n.ip, n.port), n)
		if err != nil {
			logger.Error(err)
		}
	}()
}
