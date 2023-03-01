package node

import (
	"errors"
	"github.com/madokast/GoDFS/internal/dfs/lfs"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/httputils"
	"net/http"
	"path"
)

type readReq struct {
	Path   string `json:"path,omitempty"`
	Offset int64  `json:"offset,omitempty"`
	Length int64  `json:"length,omitempty"`
}

type readRsp struct {
	Bytes []byte `json:"bytes"`
}

func (n *node) Read(path string, offset, length int64) ([]byte, error) {
	ret := web.Response[*readRsp]{}
	err := httputils.PostGob(n.ip, n.port, readApi, &readReq{Path: path, Offset: offset, Length: length}, &ret)
	if err != nil {
		return nil, err
	}
	if ret.Msg != web.SuccessMsg {
		return nil, errors.New(ret.Msg)
	}
	return ret.Data.Bytes, nil
}

func (n *node) Write(path string, offset int64, data []byte) error {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoRead(w http.ResponseWriter, r *http.Request) {
	httputils.HandleGob(w, r, &readReq{}, func(req *readReq) (*readRsp, error) {
		bytes, err := lfs.ReadLocal(path.Join(n.rootDir, req.Path), req.Offset, req.Length)
		if err != nil {
			return nil, err
		}
		return &readRsp{Bytes: bytes}, nil
	})
}

func (n *node) DoWrite(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
