package nodeimpl

import (
	"errors"
	"github.com/madokast/GoDFS/internal/dfs/dfile"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/fs"
	"github.com/madokast/GoDFS/internal/fs/lfs"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/httputils"
	"net/http"
	"path"
)

type statReq struct {
	Path string `json:"path,omitempty"`
}

type existReq struct {
	Path string `json:"path,omitempty"`
}

type existRsp struct {
	Exist bool `json:"exist,omitempty"`
}

func (n *Impl) Stat(path string) (fs.Meta, error) {
	ret := web.Response[*dfile.MetaImpl]{}
	err := httputils.PostJson(n.ip, n.port, node.StatPathApi, &statReq{Path: path}, &ret)
	if err != nil {
		return nil, err
	}
	if ret.Msg != web.SuccessMsg {
		return nil, errors.New(ret.Msg)
	}
	return ret.Data, nil
}

func (n *Impl) Exist(path string) (bool, error) {
	ret := web.Response[*existRsp]{}
	err := httputils.PostJson(n.ip, n.port, node.ExistPathApi, &existReq{Path: path}, &ret)
	if err != nil {
		return false, err
	}
	if ret.Msg != web.SuccessMsg {
		return false, errors.New(ret.Msg)
	}
	return ret.Data.Exist, nil
}

func (n *Impl) DoStat(w http.ResponseWriter, r *http.Request) {
	httputils.HandleJson(w, r, &statReq{}, func(req *statReq) (*dfile.MetaImpl, error) {
		stat, exist, err := lfs.StatLocal(path.Join(n.rootDir, req.Path))
		if err != nil {
			return nil, err
		}
		meta := &dfile.MetaImpl{
			FullName_: req.Path,
		}
		if !exist {
			meta.Exist_ = false
		} else {
			meta.Exist_ = true
			meta.Size_ = stat.Size()
			meta.IsDirectory_ = stat.IsDir()
			meta.ModifyTime_ = stat.ModTime().UnixMilli()
			meta.Locations_ = []*dfile.Location{n.Location()}
		}
		return meta, nil
	})
}

func (n *Impl) DoExist(w http.ResponseWriter, r *http.Request) {
	httputils.HandleJson(w, r, &existReq{}, func(req *existReq) (*existRsp, error) {
		existLocal := lfs.ExistLocal(path.Join(n.rootDir, req.Path))
		return &existRsp{Exist: existLocal}, nil
	})
}
