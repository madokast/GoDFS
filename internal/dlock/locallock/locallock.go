package locallock

import (
	"github.com/madokast/GoDFS/internal/dlock"
	"sync"
	"sync/atomic"
	"unsafe"
)

/**
本地化锁，用于单机模式
*/

var instance *localLock
var instancePtr unsafe.Pointer
var lock sync.Mutex

type localLock struct {
	lock sync.RWMutex
}

func New() dlock.Lock {
	pointer := atomic.LoadPointer(&instancePtr)
	if pointer == nil {
		lock.Lock()
		defer lock.Unlock()
		if atomic.LoadPointer(&instancePtr) == nil {
			instance = &localLock{}
			swapped := atomic.CompareAndSwapPointer(&instancePtr, nil, unsafe.Pointer(instance))
			if !swapped {
				panic("???")
			}
		}
	}

	return (*localLock)(atomic.LoadPointer(&instancePtr))
}

func (l *localLock) RLock() {
	l.lock.RLock()
}

func (l *localLock) RUnlock() {
	l.lock.RUnlock()
}

func (l *localLock) WLock() {
	l.lock.Lock()
}

func (l *localLock) WUnlock() {
	l.lock.Unlock()
}
