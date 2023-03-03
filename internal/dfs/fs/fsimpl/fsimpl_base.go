package fsimpl

import (
	"fmt"
	"github.com/madokast/GoDFS/internal/dfs/fs"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/dlock"
	"github.com/madokast/GoDFS/utils/logger"
	"github.com/stathat/consistent"
	"sync"
)

type Impl struct {
	allNodes        map[string]node.Node   // 素有 node
	aliveNodes      map[string]node.Node   // 存活 node
	hashCircle      *consistent.Consistent // 一致 hash
	distributedLock dlock.Lock             // 分布式锁
	localLock       sync.Mutex             // 局部锁
}

func New(distributedLock dlock.Lock, conf *fs.Conf) fs.DFS {
	dfs := &Impl{
		allNodes:        map[string]node.Node{},
		aliveNodes:      map[string]node.Node{},
		hashCircle:      consistent.New(),
		distributedLock: distributedLock,
	}
	dfs.hashCircle.NumberOfReplicas = conf.HashCircleReplicaNum
	return dfs
}

func (dfs *Impl) AddNode(n node.Node) {
	dfs.localLock.Lock()
	defer dfs.localLock.Unlock()
	dfs.allNodes[n.Key()] = n
}

func (dfs *Impl) AllNodes() []node.Node {
	dfs.localLock.Lock()
	defer dfs.localLock.Unlock()
	nodes := make([]node.Node, 0)
	for _, n := range dfs.allNodes {
		nodes = append(nodes, n)
	}
	return nodes
}

func (dfs *Impl) RefreshAliveNodesAndHashCircle() {
	dfs.localLock.Lock()
	defer dfs.localLock.Unlock()
	dfs.refreshAliveNodesAndHandCircleUnLock()
}

func (dfs *Impl) refreshAliveNodesAndHandCircleUnLock() {
	for _, n := range dfs.allNodes {
		if n.Ping() {
			_, ok := dfs.aliveNodes[n.Key()]
			if ok {
				// 已经加入过
			} else {
				logger.Info("New alive node", n.Location())
				// 需要先同步数据才能启用
				dfs.aliveNodes[n.Key()] = n
				dfs.hashCircle.Add(n.Key())
			}
		} else {
			_, ok := dfs.aliveNodes[n.Key()]
			if ok {
				// 需要移除
				logger.Warn("Lost node", n.Location())
				delete(dfs.aliveNodes, n.Key())
				dfs.hashCircle.Remove(n.Key())
			}
		}
	}
}

func (dfs *Impl) HashDistribute(key string, num int) []node.Node {
	dfs.localLock.Lock()
	defer dfs.localLock.Unlock()
	nodes, err := dfs.hashCircle.GetN(key, num)
	if err != nil {
		logger.Error(err)
	}
	ret := make([]node.Node, 0)
	for _, n := range nodes {
		ret = append(ret, dfs.allNodes[n])
	}
	return ret
}

func (dfs *Impl) String() string {
	dfs.localLock.Lock()
	defer dfs.localLock.Unlock()
	return fmt.Sprintf("DFS-Cluser(%d)", len(dfs.allNodes))
}
