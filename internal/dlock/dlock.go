package dlock

type Lock interface {
	Lock(key string)
	Unlock(key string)
}
