package fsimpl

import (
	"errors"
	"fmt"
	"github.com/madokast/GoDFS/utils/logger"
)

/**
写操作相关
*/

func (dfs *Impl) CreateFile(path string, size int64) error {
	dfs.distributedLock.WLock()
	defer dfs.distributedLock.WUnlock()
	return dfs.createFileUnlock(path, size)
}

func (dfs *Impl) createFileUnlock(path string, size int64) error {
	var errList []error
	for _, n := range dfs.HashDistribute(path, dfs.replicaNum) {
		err := n.CreateFile(path, size)
		if err != nil {
			logger.Error("DFS create", path, n.Key(), err)
			errList = append(errList, err)
		}
	}
	if len(errList) > dfs.replicaNum/2 {
		return errors.New("DFS create " + path + " " + fmt.Sprintf("%v", errList))
	}
	return nil
}

func (dfs *Impl) Delete(path string) error {
	dfs.distributedLock.WLock()
	defer dfs.distributedLock.WUnlock()
	return dfs.deleteUnlock(path)
}

func (dfs *Impl) deleteUnlock(path string) error {
	var errList []error
	for _, n := range dfs.HashDistribute(path, dfs.replicaNum) {
		err := n.Delete(path)
		if err != nil {
			logger.Error("DFS delete", path, n.Key(), err)
			errList = append(errList, err)
		}
	}
	if len(errList) > dfs.replicaNum/2 {
		return errors.New("DFS delete " + path + " " + fmt.Sprintf("%v", errList))
	}
	return nil
}

func (dfs *Impl) Write(path string, offset int64, data []byte) error {
	dfs.distributedLock.WLock()
	err := dfs.writeUnlock(path, offset, data)
	dfs.distributedLock.WUnlock()
	if err != nil {
		// refresh and try again
		dfs.refreshAliveNodesAndHandCircleUnLock()
		dfs.distributedLock.WLock()
		err = dfs.writeUnlock(path, offset, data)
		dfs.distributedLock.WUnlock()
		if err != nil {
			return err
		}
	}
	return nil
}

func (dfs *Impl) writeUnlock(path string, offset int64, data []byte) error {
	var errList []error
	for _, n := range dfs.HashDistribute(path, dfs.replicaNum) {
		err := n.Write(path, offset, data)
		if err != nil {
			logger.Error("DFS write", path, n.Key(), err)
			errList = append(errList, err)
		}
	}

	// 回调通知所有节点
	for _, n := range dfs.aliveNodes {
		err := n.WriteCallback(path, offset, int64(len(data)))
		if err != nil {
			logger.Warn("DFS write call back", path, n.String(), err)
		}
	}

	if len(errList) > dfs.replicaNum/2 {
		return errors.New("DFS write " + path + " " + fmt.Sprintf("%v", errList))
	}
	return nil
}
