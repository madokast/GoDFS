package dfs

import (
	"github.com/madokast/GoDFS/internal/dfs/fs"
	"github.com/madokast/GoDFS/internal/dfs/fs/fsimpl"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/dfs/node/nodeimpl"
	"github.com/madokast/GoDFS/internal/dlock"
)

func NewNode(conf *node.Info) node.Node {
	return nodeimpl.New(conf)
}

func NewDFS(distributedLock dlock.Lock, conf *fs.Conf) fs.DFS {
	return fsimpl.New(distributedLock, conf)
}
