package fsimpl

import (
	"errors"
	"fmt"
	"github.com/madokast/GoDFS/internal/dfs/file"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/utils/logger"
)

func (dfs *Impl) Stat(path string) (file.Meta, error) {
	dfs.distributedLock.RLock()
	defer dfs.distributedLock.RUnlock()
	return dfs.statUnlock(path)
}

func (dfs *Impl) statUnlock(path string) (file.Meta, error) {
	var errList []error
	var meta file.Meta
	for _, n := range dfs.aliveNodes {
		m, err := n.Stat(path)
		if err != nil {
			logger.Error("DFS write", path, n.Key(), err)
			errList = append(errList, err)
		} else {
			if meta == nil {
				meta = m
			}
			if m.Exist() {
				meta = m
			}
		}
	}
	if len(errList) > len(dfs.aliveNodes)/2 {
		return nil, errors.New("DFS write " + path + " " + fmt.Sprintf("%v", errList))
	}
	return meta, nil
}

func (dfs *Impl) Exist(path string) (bool, error) {
	dfs.distributedLock.RLock()
	defer dfs.distributedLock.RUnlock()
	return dfs.existUnlock(path)
}

func (dfs *Impl) existUnlock(path string) (bool, error) {
	var existNodes []node.Node
	var unExistNodes []node.Node
	var errList []error
	for _, n := range dfs.aliveNodes {
		existFile, err := n.Exist(path)
		if err != nil {
			logger.Error("DFS exist", path, n.Key(), err)
			// 出现异常当作没有
			unExistNodes = append(unExistNodes, n)
			errList = append(errList, err)
		} else {
			if existFile {
				existNodes = append(existNodes, n)
			} else {
				unExistNodes = append(unExistNodes, n)
			}
		}
	}

	var err error = nil
	if len(errList) > dfs.replicaNum/2 {
		err = errors.New("DFS exist " + path + " " + fmt.Sprintf("%v", errList))
	}

	if len(existNodes) == 0 {
		return false, err
	}

	if len(existNodes) > dfs.replicaNum/2 {
		return true, err
	} else {
		return false, err
	}

}
