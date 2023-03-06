package fsimpl

import (
	"github.com/madokast/GoDFS/internal/fs"
	"github.com/madokast/GoDFS/utils/logger"
)

func (dfs *Impl) RegisterWriteCallback(obj *fs.WriteCallBackObj) {
	if dfs.localNode == nil {
		logger.Error("No local node in DFS", dfs.String())
		return
	}
	dfs.localNode.RegisterWriteCallback(obj)
}

func (dfs *Impl) RemoveWriteCallback(obj *fs.WriteCallBackObj) {
	if dfs.localNode == nil {
		logger.Error("No local node in DFS", dfs.String())
		return
	}
	dfs.localNode.RemoveWriteCallback(obj)
}
