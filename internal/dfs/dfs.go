package dfs

import (
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/dfs/node/nodeimpl"
)

func NewNode(conf *node.Info) node.Node {
	return nodeimpl.New(conf)
}
