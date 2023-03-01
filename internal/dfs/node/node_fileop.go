package node

import (
	"github.com/madokast/GoDFS/internal/dfs"
	"net/http"
)

func (n *node) CreateFile(path string, size int64) error {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoCreateFile(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (n *node) MkdirAll(path string) {
	//TODO implement me
	panic("implement me")
}

func (n *node) ListFiles(path string) (files []string, dirs []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (n *node) Delete(path string) error {
	//TODO implement me
	panic("implement me")
}

func (n *node) Stat(path string) (dfs.FileMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (n *node) Exist(path string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoMkdirAll(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoListFiles(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoDelete(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoStat(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoExist(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
