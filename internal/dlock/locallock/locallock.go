package locallock

import (
	"github.com/madokast/GoDFS/internal/dlock"
	"github.com/madokast/GoDFS/utils/logger"
	"sync"
	"time"
)

/**
本地化锁，用于单机模式
*/

type localLock struct {
	lock sync.Mutex
	keys map[string]struct{}
}

func New() dlock.Lock {
	return &localLock{keys: map[string]struct{}{}}
}

func (l *localLock) Lock(key string) {
	start := now()
	for {
		l.lock.Lock()
		_, exist := l.keys[key]
		if exist {
			l.lock.Unlock()
			time.Sleep(10 * time.Millisecond)
			if now()-start > 10000 {
				logger.Warn(key, "deadlock?")
			}
		} else {
			l.keys[key] = struct{}{}
			l.lock.Unlock()
			break
		}
	}
}

func (l *localLock) Unlock(key string) {
	l.lock.Lock()
	defer l.lock.Unlock()
	_, exist := l.keys[key]
	if !exist {
		panic(key)
	}
	delete(l.keys, key)
}

func now() int64 {
	return time.Now().UnixMilli()
}
