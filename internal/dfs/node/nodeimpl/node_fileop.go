package nodeimpl

import (
	"errors"
	"github.com/madokast/GoDFS/internal/dfs/lfs"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/httputils"
	"net/http"
	"path"
)

type createFileReq struct {
	Path string `json:"path"`
	Size int64  `json:"size"`
}

type deleteFileReq struct {
	Path string `json:"path"`
}

func (n *Impl) CreateFile(path string, size int64) error {
	ret := web.Response[*web.NullResponse]{}
	err := httputils.PostJson(n.ip, n.port, createFileApi, &createFileReq{Path: path, Size: size}, &ret)
	if err != nil {
		return err
	}
	if ret.Msg != web.SuccessMsg {
		return errors.New(ret.Msg)
	}
	return nil
}

func (n *Impl) DoCreateFile(w http.ResponseWriter, r *http.Request) {
	httputils.HandleJson(w, r, &createFileReq{}, func(req *createFileReq) (*web.NullResponse, error) {
		err := lfs.CreateFileLocal(path.Join(n.rootDir, req.Path), req.Size)
		return &web.NullResponse{}, err
	})
}

// Delete 删除文件、文件夹，递归删。不存在不报错
func (n *Impl) Delete(path string) error {
	ret := web.Response[*web.NullResponse]{}
	err := httputils.PostJson(n.ip, n.port, deletePathApi, &deleteFileReq{Path: path}, &ret)
	if err != nil {
		return err
	}
	if ret.Msg != web.SuccessMsg {
		return errors.New(ret.Msg)
	}
	return nil
}

func (n *Impl) DoDelete(w http.ResponseWriter, r *http.Request) {
	httputils.HandleJson(w, r, &deleteFileReq{}, func(req *deleteFileReq) (*web.NullResponse, error) {
		err := lfs.DeleteLocal(path.Join(n.rootDir, req.Path))
		return &web.NullResponse{}, err
	})
}
