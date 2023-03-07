package lru

import (
	"container/list"
	"github.com/madokast/GoDFS/internal/allocator"
	"github.com/madokast/GoDFS/internal/fs/write_callback"
	"sync"
)

/**
lru 缓存
*/

type node struct {
	key   allocator.Pointer
	value *allocator.CacheData
}

// 每个缓存额外大小：Pointer 8 + CacheData.func 8 + *node 8 + element 8
const extraSz = 32

type CacheLRU struct {
	maxSize uint64
	curSize uint64
	wcb     write_callback.Register
	list    *list.List
	cache   map[allocator.Pointer]*list.Element
	lock    sync.RWMutex
}

func New(wcb write_callback.Register, maxSize uint64) *CacheLRU {
	return &CacheLRU{
		maxSize: maxSize,
		wcb:     wcb,
		list:    list.New(),
		cache:   map[allocator.Pointer]*list.Element{},
	}
}

func (lru *CacheLRU) Put(p allocator.Pointer, data *allocator.CacheData) {
	//logger.Debug("Cache put", &p, data.Data)
	// 注册缓存
	lru.wcb.RegisterWriteCallback(data.WcObj)

	lru.lock.Lock()
	defer lru.lock.Unlock()
	ele, ok := lru.cache[p]
	if ok {
		lru.list.MoveToFront(ele)
		return
	}
	ele = lru.list.PushFront(&node{
		key:   p,
		value: data,
	})
	lru.cache[p] = ele
	lru.curSize += uint64(len(data.Data)) + extraSz

	lru.expire()
}

func (lru *CacheLRU) Get(p allocator.Pointer) (data *allocator.CacheData, ok bool) {
	lru.lock.RLock()
	defer lru.lock.RUnlock()
	ele, ok := lru.cache[p]
	if ok {
		lru.list.MoveToFront(ele)
		return ele.Value.(*node).value, true
	}
	return nil, false
}

// Remove 只会被回调函数或者 expire 调用
func (lru *CacheLRU) Remove(p allocator.Pointer) {
	//logger.Debug("Cache remove", &p)
	lru.lock.Lock()
	defer lru.lock.Unlock()
	ele, ok := lru.cache[p]
	if ok {
		delete(lru.cache, ele.Value.(*node).key)
		lru.list.Remove(ele)
		go lru.wcb.RemoveWriteCallback(ele.Value.(*node).value.WcObj)
	}
}

func (lru *CacheLRU) expire() {
	for lru.curSize > lru.maxSize && lru.list.Len() > 0 {
		back := lru.list.Back()
		// 新大小
		lru.curSize -= uint64(len(back.Value.(*node).value.Data)) + extraSz
		// 移除回调
		lru.wcb.RemoveWriteCallback(back.Value.(*node).value.WcObj)
		// 删除 map
		delete(lru.cache, back.Value.(*node).key)
		// 移除 list
		lru.list.Remove(back)
	}
}
