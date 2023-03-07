package lru

import (
	"fmt"
	"github.com/madokast/GoDFS/internal/allocator"
	"github.com/madokast/GoDFS/internal/fs/write_callback"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/logger"
	"net/http"
	"testing"
)

func TestCacheLRU_Put(t *testing.T) {
	lru := New(&wcb{}, 1024)
	for i := 0; i < 20; i++ {
		fileName := fmt.Sprintf("F%d.dat", i)
		logger.Info("Put", fileName)
		lru.Put(allocator.Pointer{BlockId: uint32(i), BlockOffset: uint32(i + i)},
			&allocator.CacheData{Data: make([]byte, 150),
				WcObj: &write_callback.Entry{
					FileName: fileName,
					Callback: func() {
						logger.Info("remove", fileName)
					},
				}})
	}
}

func TestCacheLRU_Get(t *testing.T) {
	lru := New(&wcb{}, 1024)
	for i := 0; i < 20; i++ {
		fileName := fmt.Sprintf("F%d.dat", i)
		logger.Info("Put", fileName)
		pointer := allocator.Pointer{BlockId: uint32(i), BlockOffset: uint32(i + i)}
		lru.Put(pointer,
			&allocator.CacheData{Data: make([]byte, 150),
				WcObj: &write_callback.Entry{
					FileName: fileName,
					Callback: func() {
						logger.Info("remove", fileName)
					},
				}})
		data, ok := lru.Get(pointer)
		utils.PanicIf(!ok)
		utils.PanicIf(data.WcObj.FileName != fileName)
	}
}

type wcb struct {
}

func (wc *wcb) RegisterWriteCallback(*write_callback.Entry) {
	return
}

func (wc *wcb) RemoveWriteCallback(obj *write_callback.Entry) {
	logger.Info("Remove", obj.FileName)
	return
}

func (wc *wcb) WriteCallback(string, int64, int64) error {
	return nil
}

func (wc *wcb) DoWriteCallback(http.ResponseWriter, *http.Request) {
	return
}
