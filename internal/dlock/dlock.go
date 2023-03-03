package dlock

type Lock interface {
	RLock()
	RUnlock()
	WLock()
	WUnlock()
}
