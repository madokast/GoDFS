package nodeimpl

import (
	"errors"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"net/http"
)

type syncReq struct {
	SrcNode *node.Info `json:"srcNode"`
	File    string     `json:"file"`
}

func (n *Impl) Sync(src node.Node, file string) error {
	if n.Key() == src.Key() {
		return nil
	}

	ret := web.Response[*web.NullResponse]{}
	err := httputils.PostJson(n.ip, n.port, node.SyncApi, &syncReq{
		SrcNode: src.Info(),
		File:    file,
	}, &ret)
	if err != nil {
		return err
	}
	if ret.Msg != web.SuccessMsg {
		return errors.New(ret.Msg)
	}
	return nil
}

func (n *Impl) DoSync(w http.ResponseWriter, r *http.Request) {
	httputils.HandleJson(w, r, &syncReq{}, func(req *syncReq) (*web.NullResponse, error) {
		src := New(req.SrcNode)

		// 我方是否存在
		exist, err := n.Exist(req.File)
		if err != nil {
			return nil, err
		}
		if exist {
			// 存在，判断 MD5
			thisMd5, err := n.MD5(req.File)
			if err != nil {
				return nil, err
			}
			thatMd5, err := src.MD5(req.File)
			if err != nil {
				return nil, err
			}
			if thisMd5 == thatMd5 {
				logger.Debug("Same MD5", thisMd5, req.File, n.Key(), src.Key(), "sync pass")
				return nil, nil
			} else {
				// 复制
				logger.Debug("Inconsistent MD5", thisMd5, thatMd5, req.File, n.Key(), src.Key())
				goto copy
			}
		} else {
			// 不存在创建
			stat, err := src.Stat(req.File)
			if err != nil {
				return nil, err
			}
			if !stat.Exist() {
				return nil, errors.New("Src node " + src.Key() + " contains no " + req.File)
			}
			err = n.CreateFile(req.File, stat.Size())
			if err != nil {
				return nil, err
			}
			// 复制
			goto copy
		}
	copy:
		stat, err := n.Stat(req.File)
		if err != nil {
			return nil, err
		}
		if !stat.Exist() {
			return nil, errors.New("The node " + n.Key() + " " + req.File + " concurrent modified?")
		}
		read, err := src.Read(req.File, 0, stat.Size())
		if err != nil {
			return nil, err
		}
		err = n.Write(req.File, 0, read)
		if err != nil {
			return nil, err
		}
		logger.Info("Sync done", req.File, src.Key(), "to", n.Key())
		return nil, nil
	})
}
