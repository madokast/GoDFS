package fsimpl

import "github.com/madokast/GoDFS/internal/dfs/file"

func (dfs *Impl) CreateFile(path string, size int64) error {
	//TODO implement me
	panic("implement me")
}

func (dfs *Impl) ListFiles(path string) (files []string, dirs []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (dfs *Impl) Delete(path string) error {
	//TODO implement me
	panic("implement me")
}

func (dfs *Impl) Stat(path string) (file.Meta, error) {
	//TODO implement me
	panic("implement me")
}

func (dfs *Impl) Exist(path string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (dfs *Impl) Read(path string, offset, length int64) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (dfs *Impl) Write(path string, offset int64, data []byte) error {
	//TODO implement me
	panic("implement me")
}
