package allocator_impl

import (
	"errors"
	"fmt"
	"github.com/madokast/GoDFS/internal/allocator"
	"github.com/madokast/GoDFS/internal/allocator/lru"
	fs2 "github.com/madokast/GoDFS/internal/fs"
	"github.com/madokast/GoDFS/internal/fs/write_callback"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/logger"
	"path"
	"sync"
)

type impl struct {
	fs            fs2.BaseFS
	lock          sync.Mutex    // 锁
	lru           *lru.CacheLRU // LRU 缓存
	baseDir       string        // block 文件所在目录
	blockFileSize uint32        // 一个 block file 的大小
}

// 第一个 block file 的 0 位置存储了当前 FREE 指针值
var firstBlockFile = blockFileName(0)

func New(fs fs2.BaseFS, wcr write_callback.Register, baseDir string, blockFileSize uint32, cacheMaxSize uint64) (allocator.Allocator, error) {
	fs.WLock()
	defer fs.WUnlock()

	files, dirs, err := fs.ListFilesUnlock(baseDir)
	if err != nil {
		return nil, errors.New("New allocator failed because " + err.Error())
	}
	if len(files) > 0 || len(dirs) > 0 {
		logger.Warn("New allocator base dir " + baseDir + " is not empty")
		err = fs.DeleteUnlock(baseDir)
		if err != nil {
			return nil, errors.New("New allocator failed when clear base dir because " + err.Error())
		}
	}
	firstBlockPath := path.Join(baseDir, firstBlockFile)
	err = fs.CreateFileUnlock(firstBlockPath, int64(blockFileSize))
	if err != nil {
		return nil, errors.New("New allocator failed when creating first block file because " + err.Error())
	}
	// 因为 (0,0) 位置存放 firstFreePointer 的值
	// 所以 firstFreePointer 值为 (0,4)
	firstFreePointer := allocator.Pointer{BlockId: 0, BlockOffset: allocator.PointerSz}
	err = fs.WriteUnlock(firstBlockPath, 0, firstFreePointer.ToBytes())
	if err != nil {
		return nil, errors.New("New allocator failed when store first free pointer because " + err.Error())
	}

	alloc := &impl{
		fs:            fs,
		lru:           lru.New(wcr, cacheMaxSize),
		baseDir:       baseDir,
		blockFileSize: blockFileSize,
	}
	logger.Info("Create allocator " + alloc.String())
	return alloc, nil
}

// Allocate 分配内存，注意分布式全局锁
// newBlock 强制在新 block 文件上分配内存
func (i *impl) Allocate(size uint32, newBlock bool) (allocator.Pointer, error) {
	// size 加上头信息
	sizeWithHeader := size + MemHeaderSz
	firstBlockPath := path.Join(i.baseDir, firstBlockFile)

	i.fs.WLock()
	defer i.fs.WUnlock()

	ptrBytes, err := i.fs.ReadUnlock(firstBlockPath, 0, int64(allocator.PointerSz))
	if err != nil {
		return allocator.Pointer{}, errors.New("Alloc Error in read first block file because " + err.Error())
	}
	nextPtr := allocator.ReadPointer(ptrBytes)
	//logger.Debug("nextPtr", nextPtr)

	// 判断是不是分配到新文件
	if nextPtr.BlockOffset+sizeWithHeader > i.blockFileSize {
		newBlock = true
	}

	// 返回的地址
	var retPtr allocator.Pointer

	// 新文件
	if newBlock {
		newBlockId := nextPtr.BlockId + 1
		newBlockFile := path.Join(i.baseDir, blockFileName(newBlockId))
		err = i.fs.CreateFileUnlock(newBlockFile, int64(i.blockFileSize))
		if err != nil {
			err2 := i.fs.DeleteUnlock(newBlockFile)
			if err2 != nil {
				logger.Error(err2)
			}
			return allocator.Pointer{}, errors.New("Alloc Error in create new block file because " + err.Error())
		}
		retPtr = allocator.Pointer{BlockId: newBlockId, BlockOffset: 0}
		nextPtr = allocator.Pointer{BlockId: newBlockId, BlockOffset: sizeWithHeader}
	} else {
		// 旧文件
		retPtr = nextPtr
		nextPtr.BlockOffset += sizeWithHeader
	}

	// nextPtr 写回
	err = i.fs.WriteUnlock(firstBlockPath, 0, nextPtr.ToBytes())
	if err != nil {
		return allocator.Pointer{}, errors.New("Alloc Error in write next pointer because " + err.Error())
	}

	// 从 retPtr 读取信息，头部写入 size 信息
	blockPath := path.Join(i.baseDir, blockFileName(retPtr.BlockId))
	read, err := i.fs.ReadUnlock(blockPath, int64(retPtr.BlockOffset), int64(sizeWithHeader))
	if err != nil {
		return allocator.Pointer{}, errors.New("Alloc Error in read allocated space because " + err.Error())
	}
	header := CreateHeader(size)
	header.WriteHeaderTo(read)

	// read 写回
	err = i.fs.WriteUnlock(blockPath, int64(retPtr.BlockOffset), read[:MemHeaderSz])
	if err != nil {
		return allocator.Pointer{}, errors.New("Alloc Error in write header because " + err.Error())
	}

	// 注册写监听，缓存。注意存储的数据为 read[MemHeaderSz:] 不要头
	cache := &allocator.CacheData{
		Data: read[MemHeaderSz:], // 数据
		WcObj: &write_callback.Entry{
			FileName: blockPath,                                                                              // 文件名
			Offset:   int64(retPtr.BlockOffset),                                                              // 偏移
			Length:   int64(sizeWithHeader),                                                                  // 长度带有头
			Callback: func(catch allocator.Pointer) func() { return func() { i.lru.Remove(catch) } }(retPtr), // 回调移除
		},
	}
	// 缓存
	i.lru.Put(retPtr, cache)

	return retPtr, nil
}

// Free 释放内存。仅仅是标记，以后实现物理释放。
func (i *impl) Free(p allocator.Pointer) error {
	i.fs.WLock()
	defer i.fs.WUnlock()

	// 读取 header
	blockPath := path.Join(i.baseDir, blockFileName(p.BlockId))
	headerBytes, err := i.fs.ReadUnlock(blockPath, int64(p.BlockOffset), int64(MemHeaderSz))
	if err != nil {
		return errors.New("Free Error in read header because " + err.Error())
	}
	header := ReadHeader(headerBytes)

	// free
	header.Free()
	header.WriteHeaderTo(headerBytes)

	// 写回
	err = i.fs.Write(blockPath, int64(p.BlockOffset), headerBytes)
	if err != nil {
		return errors.New("Free Error in write header because " + err.Error())
	}

	// 删除缓存
	i.lru.Remove(p)
	return nil
}

func (i *impl) Read(p allocator.Pointer, reader func([]byte)) error {
	// 尝试缓存中读取
	data, ok := i.lru.Get(p)
	if ok {
		reader(data.Data)
		return nil
	} else {
		return i.readFromDFS(p, reader)
	}
}

func (i *impl) readFromDFS(p allocator.Pointer, reader func([]byte)) error {
	i.fs.RLock()
	defer i.fs.RUnlock()

	// 从 DFS 中读取，直接读 1 MB
	blockPath := path.Join(i.baseDir, blockFileName(p.BlockId))
	bytes, err := i.fs.ReadUnlock(blockPath, int64(p.BlockOffset), min(1024*1024, int64(i.blockFileSize-p.BlockOffset)))
	if err != nil {
		return errors.New("Read pointer error in read dfs " + blockPath + " because " + err.Error())
	}
	// 1MB 解析
	header := ReadHeader(bytes)
	if header.IsFree() {
		panic("Read a released pointer " + p.String())
	}
	size := header.Size()
	if size > uint32(len(bytes))-MemHeaderSz {
		// 读少了，读取剩余没有读的数目
		// 已经从 offset 读取了 len(bytes) 数据，实际数据为 len(bytes)-MemHeaderSz
		// 发现数据长度为 size，那么还要读 size - (len(bytes)-MemHeaderSz) 的数据
		// 偏移为 offset + len(bytes)
		bytes2, err := i.fs.ReadUnlock(blockPath, int64(p.BlockOffset+uint32(len(bytes))), int64(size-(uint32(len(bytes))-MemHeaderSz)))
		if err != nil {
			return errors.New("Read pointer error in read dfs " + blockPath + " because " + err.Error())
		}
		data := append(bytes[MemHeaderSz:], bytes2...)
		utils.PanicIf(uint32(len(data)) != size, p, len(data), size)
		// 加入缓存
		i.lru.Put(p, &allocator.CacheData{
			Data: data,
			WcObj: &write_callback.Entry{
				FileName: blockPath,
				Offset:   int64(p.BlockOffset),
				Length:   int64(size + MemHeaderSz),
				Callback: func(catch allocator.Pointer) func() { return func() { i.lru.Remove(catch) } }(p),
			},
		})
		// 回调
		reader(data)
	} else {
		// 读多了，先解析出第一个指针数据
		data := bytes[MemHeaderSz : size+MemHeaderSz]
		// 加入缓存
		i.lru.Put(p, &allocator.CacheData{
			Data: data,
			WcObj: &write_callback.Entry{
				FileName: blockPath,
				Offset:   int64(p.BlockOffset),
				Length:   int64(size + MemHeaderSz),
				Callback: func(catch allocator.Pointer) func() { return func() { i.lru.Remove(catch) } }(p),
			},
		})
		// 回调
		reader(data)

		// 剩余数据看是否存在完整的指针数据，
		//remainBytes := bytes[size+MemHeaderSz:]
		//nextPtr := allocator.Pointer{BlockId: p.BlockId, BlockOffset: p.BlockOffset + size + MemHeaderSz}
		//for {
		//	header := ReadHeader(remainBytes)
		//	if header.Size() <= uint32(len(remainBytes))-MemHeaderSz {
		//		data := remainBytes[MemHeaderSz : header.Size()+MemHeaderSz]
		//		i.lru.Put(nextPtr, &allocator.CacheData{
		//			Data: data,
		//			WcObj: &write_callback.Entry{
		//				FileName: blockPath,
		//				Offset:   int64(nextPtr.BlockOffset),
		//				Length:   int64(header.Size() + MemHeaderSz),
		//				Callback: func(catch allocator.Pointer) func() { return func() { i.lru.Remove(catch) } }(nextPtr),
		//			},
		//		})
		//	} else {
		//		break
		//	}
		//	remainBytes = remainBytes[header.Size()+MemHeaderSz:]
		//	nextPtr = allocator.Pointer{BlockId: p.BlockId, BlockOffset: nextPtr.BlockOffset + header.Size() + MemHeaderSz}
		//}
	}
	return nil
}

func (i *impl) ReadString(p allocator.Pointer) (string, error) {
	var s string
	err := i.Read(p, func(bytes []byte) {
		s = string(bytes)
	})
	return s, err
}

func (i *impl) ReadBytes(p allocator.Pointer) ([]byte, error) {
	var b []byte
	err := i.Read(p, func(bytes []byte) {
		b = make([]byte, len(bytes))
		copy(b, bytes)
	})
	return b, err
}

func (i *impl) Write(p allocator.Pointer, offset uint32, data []byte) error {
	blockPath := path.Join(i.baseDir, blockFileName(p.BlockId))
	i.fs.WLock()
	defer i.fs.WUnlock()
	// 不需要失效缓存，会自动通过回调移除
	return i.fs.WriteUnlock(blockPath, int64(p.BlockOffset+MemHeaderSz+offset), data)
}

func (i *impl) String() string {
	return fmt.Sprintf("Alloc[%s](%s)", i.baseDir, i.fs.String())
}

func blockFileName(bid uint32) string {
	return fmt.Sprintf("%010d.dat", bid)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}
