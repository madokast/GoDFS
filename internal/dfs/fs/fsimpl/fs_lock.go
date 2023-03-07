package fsimpl

func (dfs *Impl) LocalLock() {
	dfs.localLock.Lock()
}

func (dfs *Impl) LocalUnlock() {
	dfs.localLock.Unlock()
}

func (dfs *Impl) DistributedRLock() {
	dfs.distributedLock.RLock()
}

func (dfs *Impl) DistributedRUnlock() {
	dfs.distributedLock.RUnlock()
}

func (dfs *Impl) DistributedWLock() {
	dfs.distributedLock.WLock()
}

func (dfs *Impl) DistributedWUnlock() {
	dfs.distributedLock.WUnlock()
}

func (dfs *Impl) RLock() {
	dfs.DistributedRLock()
}

func (dfs *Impl) RUnlock() {
	dfs.DistributedRUnlock()
}

func (dfs *Impl) WLock() {
	dfs.DistributedWLock()
}

func (dfs *Impl) WUnlock() {
	dfs.DistributedWUnlock()
}
