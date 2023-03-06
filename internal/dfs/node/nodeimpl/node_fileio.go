package nodeimpl

import (
	"errors"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/fs/lfs"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/httputils"
	"net/http"
	"path"
)

type readReq struct {
	Path   string
	Offset int64
	Length int64
}

type readRsp struct {
	Bytes []byte
}

type writeReq struct {
	Path   string
	Offset int64
	Bytes  []byte
}

func (n *Impl) Read(path string, offset, length int64) ([]byte, error) {
	ret := web.Response[*readRsp]{}
	err := httputils.PostGob(n.ip, n.port, node.ReadFileApi, &readReq{Path: path, Offset: offset, Length: length}, &ret)
	if err != nil {
		return nil, err
	}
	if ret.Msg != web.SuccessMsg {
		return nil, errors.New(ret.Msg)
	}
	return ret.Data.Bytes, nil
}

func (n *Impl) Write(path string, offset int64, data []byte) error {
	ret := web.Response[*web.NullResponse]{}
	err := httputils.PostGob(n.ip, n.port, node.WriteFileApi, &writeReq{Path: path, Offset: offset, Bytes: data}, &ret)
	if err != nil {
		return err
	}
	if ret.Msg != web.SuccessMsg {
		return errors.New(ret.Msg)
	}
	return nil
}

func (n *Impl) DoRead(w http.ResponseWriter, r *http.Request) {
	httputils.HandleGob(w, r, &readReq{}, func(req *readReq) (*readRsp, error) {
		bytes, err := lfs.ReadLocal(path.Join(n.rootDir, req.Path), req.Offset, req.Length)
		if err != nil {
			return nil, err
		}
		return &readRsp{Bytes: bytes}, nil
	})
}

func (n *Impl) DoWrite(w http.ResponseWriter, r *http.Request) {
	httputils.HandleGob(w, r, &writeReq{}, func(req *writeReq) (*web.NullResponse, error) {
		err := lfs.WriteLocal(path.Join(n.rootDir, req.Path), req.Offset, req.Bytes)
		return &web.NullResponse{}, err
	})
}
