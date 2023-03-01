package node

import "net/http"

func (n *node) Lock(path string) error {
	//TODO implement me
	panic("implement me")
}

func (n *node) Unlock(path string) error {
	//TODO implement me
	panic("implement me")
}

func (n *node) SetVersion(path string) error {
	//TODO implement me
	panic("implement me")
}

func (n *node) Version(path string) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (n *node) Read(path string, offset, length int64) ([]byte, error) {
	//TODO implement me
	panic("implement me")
}

func (n *node) Write(path string, offset int64, data []byte) error {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoLock(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoUnlock(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoSetVersion(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoVersion(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoRead(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}

func (n *node) DoWrite(w http.ResponseWriter, r *http.Request) {
	//TODO implement me
	panic("implement me")
}
