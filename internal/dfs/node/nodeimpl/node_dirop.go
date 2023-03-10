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

type listFilesReq struct {
	Path string `json:"path,omitempty"`
}

type listFilesRsp struct {
	Files []string `json:"files,omitempty"`
	Dirs  []string `json:"dirs,omitempty"`
}

func (n *Impl) ListFiles(path string) (files []string, dirs []string, err error) {
	ret := web.Response[*listFilesRsp]{}
	err = httputils.PostJson(n.ip, n.port, node.ListFilesApi, &listFilesReq{Path: path}, &ret)
	if err != nil {
		return nil, nil, err
	}
	if ret.Msg != web.SuccessMsg {
		return nil, nil, errors.New(ret.Msg)
	}
	return ret.Data.Files, ret.Data.Dirs, nil
}

func (n *Impl) DoListFiles(w http.ResponseWriter, r *http.Request) {
	httputils.HandleJson(w, r, &listFilesReq{}, func(req *listFilesReq) (*listFilesRsp, error) {
		files, dirs, err := lfs.ListFilesLocal(path.Join(n.rootDir, req.Path))
		for i := 0; i < len(files); i++ {
			files[i] = path.Join(req.Path, files[i])
		}
		for i := 0; i < len(dirs); i++ {
			dirs[i] = path.Join(req.Path, dirs[i])
		}
		return &listFilesRsp{Files: files, Dirs: dirs}, err
	})
}

/*================ 有锁无锁相同 ======================*/

func (n *Impl) ListFilesUnlock(path string) (files []string, dirs []string, err error) {
	return n.ListFiles(path)
}
