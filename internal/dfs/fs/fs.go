package fs

import (
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/fs"
)

/**
分布式文件系统
*/

type Conf struct {
	HashCircleReplicaNum int
	FileReplicaNum       int
}

type DFS interface {
	fs.BaseFS
	dfsBase
	writeCallback
}

type dfsBase interface {
	AddNode(node node.Node)
	AllNodes() []node.Node
	RefreshAliveNodesAndHashCircle() // 刷新存活的 node 和 hash 环，定时调用
	HashDistribute(key string, num int) []node.Node
	String() string
}

type writeCallback interface {
	RegisterWriteCallback(*fs.WriteCallBackObj) // 注册文件修改通知回调。缓存层需要用到，用来失效一些资源
	RemoveWriteCallback(*fs.WriteCallBackObj)   // 取消注册
}
