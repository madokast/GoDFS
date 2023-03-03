package fsimpl

import (
	"errors"
	"fmt"
	"github.com/madokast/GoDFS/utils/logger"
)

func (dfs *Impl) ListFiles(path string) (files []string, dirs []string, err error) {
	var errList []error
	allFiles := map[string]struct{}{}
	allDirs := map[string]struct{}{}
	for _, n := range dfs.aliveNodes {
		files, dirs, err = n.ListFiles(path)
		if err != nil {
			logger.Error("DFS list", path, n.Key(), err)
			errList = append(errList, err)
		} else {
			for _, file := range files {
				allFiles[file] = struct{}{}
			}
			for _, dir := range dirs {
				allDirs[dir] = struct{}{}
			}
		}
	}
	if len(errList) > dfs.replicaNum/2 {
		return nil, nil, errors.New("DFS write " + path + " " + fmt.Sprintf("%v", errList))
	}

	files = make([]string, 0, len(allFiles))
	dirs = make([]string, 0, len(allDirs))
	for file := range allFiles {
		files = append(files, file)
	}
	for dir := range allDirs {
		dirs = append(dirs, dir)
	}

	return files, dirs, nil
}

func (dfs *Impl) Read(path string, offset, length int64) ([]byte, error) {
	var errList []error
	var read []byte
	for _, n := range dfs.HashDistribute(path, dfs.replicaNum) {
		bytes, err := n.Read(path, offset, length)
		if err != nil {
			logger.Error("DFS write", path, n.Key(), err)
			errList = append(errList, err)
		} else {
			read = bytes
			break
		}
	}
	if len(errList) > dfs.replicaNum/2 {
		return read, errors.New("DFS write " + path + " " + fmt.Sprintf("%v", errList))
	}
	return read, nil
}
