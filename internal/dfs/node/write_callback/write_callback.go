package write_callback

import (
	"errors"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/fs/write_callback"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/httputils"
	"net/http"
	"sync"
)

/**
不要直接使用这个类，这个类是 node 的内部类
*/

type objSet = map[*write_callback.Entry]struct{}

type Impl struct {
	sync.Mutex
	callbackMap map[string]objSet
	node        node.Node
}

func New(node node.Node) write_callback.WriteCallBack {
	return &Impl{callbackMap: map[string]objSet{}, node: node}
}

func (wc *Impl) RegisterWriteCallback(obj *write_callback.Entry) {
	wc.Lock()
	defer wc.Unlock()
	set, ok := wc.callbackMap[obj.FileName]

	if !ok {
		set = objSet{}
	}

	set[obj] = struct{}{}
	wc.callbackMap[obj.FileName] = set
}

func (wc *Impl) RemoveWriteCallback(obj *write_callback.Entry) {
	wc.Lock()
	defer wc.Unlock()
	wc.removeWriteCallbackUnlock(obj)
}

func (wc *Impl) removeWriteCallbackUnlock(obj *write_callback.Entry) {
	set, ok := wc.callbackMap[obj.FileName]
	if ok {
		delete(set, obj)
		wc.callbackMap[obj.FileName] = set
	}
}

type callbackReq struct {
	Path   string
	Offset int64
	Length int64
}

func (wc *Impl) WriteCallback(fileName string, offset, length int64) error {
	ret := web.Response[*web.NullResponse]{}
	err := httputils.PostJson(wc.node.IP(), wc.node.Port(), node.WriteCallBackApi,
		&callbackReq{Path: fileName, Offset: offset, Length: length}, &ret)
	if err != nil {
		return err
	}
	if ret.Msg != web.SuccessMsg {
		return errors.New(ret.Msg)
	}
	return nil
}

func (wc *Impl) DoWriteCallback(w http.ResponseWriter, r *http.Request) {
	httputils.HandleJson(w, r, &callbackReq{}, func(req *callbackReq) (*web.NullResponse, error) {
		wc.Lock()
		defer wc.Unlock()
		//logger.Debug("Write callback", req.Path, req.Offset, req.Length)
		set, ok := wc.callbackMap[req.Path]
		if ok {
			var called []*write_callback.Entry
			for obj := range set {
				if obj.Intersect(req.Offset, req.Length) {
					//logger.Debug("callback in", wc.node.String())
					obj.Callback()
					called = append(called, obj)
				}
			}

			// remove
			for _, obj := range called {
				wc.removeWriteCallbackUnlock(obj)
			}
		}
		return &web.NullResponse{}, nil
	})
}
