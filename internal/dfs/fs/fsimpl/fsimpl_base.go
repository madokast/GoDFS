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
	distributedLock dlock.Lock             // 分布式锁，对文件/目录加锁。不管目录下的文件
	localLock       sync.Mutex             // 局部锁
	replicaNum      int                    // 副本数目
}

func New(distributedLock dlock.Lock, conf *fs.Conf) fs.DFS {
	conf = checkConf(conf)
	dfs := &Impl{
		allNodes:        map[string]node.Node{},
		aliveNodes:      map[string]node.Node{},
		hashCircle:      consistent.New(),
		distributedLock: distributedLock,
		replicaNum:      conf.FileReplicaNum,
	}
	dfs.hashCircle.NumberOfReplicas = conf.HashCircleReplicaNum
	return dfs
}

func checkConf(conf *fs.Conf) *fs.Conf {
	if conf == nil {
		return &fs.Conf{
			HashCircleReplicaNum: 10,
			FileReplicaNum:       3,
		}
	}
	if conf.HashCircleReplicaNum <= 0 {
		conf.HashCircleReplicaNum = 10
	}
	if conf.FileReplicaNum <= 0 {
		conf.HashCircleReplicaNum = 3
	}
	return conf
}

// AddNode 添加节点，仅仅加入 allNodes 中
func (dfs *Impl) AddNode(n node.Node) {
	dfs.localLock.Lock()
	defer dfs.localLock.Unlock()
	dfs.allNodes[n.Key()] = n
}

// AllNodes 返回 allNodes
func (dfs *Impl) AllNodes() []node.Node {
	dfs.localLock.Lock()
	defer dfs.localLock.Unlock()
	nodes := make([]node.Node, 0)
	for _, n := range dfs.allNodes {
		nodes = append(nodes, n)
	}
	return nodes
}

// RefreshAliveNodesAndHashCircle 刷线所有存活的 node，同时刷新 hash 环
// 如果遇到存活节点的变化，需要重新分布数据
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
				// 新节点加入，newNode 数据需要分配
				logger.Info("New alive node", n.Location())
				dfs.distributedLock.WLock()
				dfs.aliveNodes[n.Key()] = n
				dfs.hashCircle.Add(n.Key())
				dfs.nodeSyncUnlock(n)
				dfs.distributedLock.WUnlock()
				logger.Info("Finish new node adding", n.Location())
			}
		} else {
			_, ok := dfs.aliveNodes[n.Key()]
			if ok {
				// 需要移除
				logger.Warn("Lost node", n.Location())
				dfs.distributedLock.WLock()
				delete(dfs.aliveNodes, n.Key())
				dfs.hashCircle.Remove(n.Key())
				for _, eachNode := range dfs.aliveNodes {
					dfs.nodeSyncUnlock(eachNode)
				}
				dfs.distributedLock.WUnlock()
			}
		}
	}
}

func (dfs *Impl) HashDistribute(key string, num int) []node.Node {
	dfs.localLock.Lock()
	defer dfs.localLock.Unlock()
	return dfs.hashDistributeUnlock(key, num)
}

func (dfs *Impl) hashDistributeUnlock(key string, num int) []node.Node {
	nodes, err := dfs.hashCircle.GetN(key, num)
	if err != nil {
		// 几乎不可能，因为自己就有一个 node
		logger.Error("No available nodes?", dfs.AllNodes(), err)
		return nil
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

// nodeSyncUnlock 节点重新同步数据
// 1. 同步其他节点的数据（会 hash 到自己的那些文件）
// 2. 自身数据清理
func (dfs *Impl) nodeSyncUnlock(newNode node.Node) {
	// 同步其他节点的数据
	for _, aliveNode := range dfs.aliveNodes {
		dfs.syncUnlock(newNode, aliveNode, "/")
	}

	// 自身数据清理，删除那些不再属于自己的数据
	newNode.ForAllFile("/", func(file string) {
		remove := true
		for _, n := range dfs.hashDistributeUnlock(file, dfs.replicaNum) {
			if n.Key() == newNode.Key() {
				remove = false // 应该继续保存该文件
			}
		}
		if remove {
			logger.Debug("DFS reduce replica", file, "in", newNode.Key())
			err := newNode.Delete(file)
			if err != nil {
				logger.Error("DFS reduce replica", file, "in", newNode.Key(), err)
			}
		}
	})
}

func (dfs *Impl) sync(des node.Node, src node.Node, path string) {
	dfs.distributedLock.WLock()
	defer dfs.distributedLock.WUnlock()
	dfs.syncUnlock(des, src, path)
}

// syncUnlock 将 src 中 path 目录下全部文件（遍历）同步到 des 中。
// 只同步那些需要到 des 中的文件
// 同时，src 中的文件，如果不在 hash 环最近 replicaNum 了，也会删除
func (dfs *Impl) syncUnlock(des node.Node, src node.Node, path string) {
	if des.Key() == src.Key() {
		return
	}

	src.ForAllFile(path, func(file string) {
		replicaNodes := dfs.hashDistributeUnlock(file, dfs.replicaNum)
		srcHolderFile := false // src node 是否应该继续保存改文件
		for _, n := range replicaNodes {
			if n.Key() == des.Key() {
				logger.Debug("DFS sync", file, "from", src.Key(), "to", des.Key())
				err := des.Sync(src, file)
				if err != nil {
					logger.Error("DFS sync", file, "from", src.Key(), "to", des.Key(), err)
				}
			}
			// 如果 src 仍在最近的 replicaNum 中，就不删除
			if n.Key() == src.Key() {
				srcHolderFile = true
			}
		}
		// src 不应该保存文件 file 了
		if !srcHolderFile {
			logger.Debug("DFS sync", file, "reduce replica", src.Key())
			err := src.Delete(file)
			if err != nil {
				logger.Error("DFS sync", file, "reduce replica", src.Key(), err)
			}
		}
	})
}
