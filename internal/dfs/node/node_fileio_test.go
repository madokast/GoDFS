package node

import (
	"github.com/madokast/GoDFS/internal/dfs"
	"github.com/madokast/GoDFS/internal/dfs/lfs"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
	"time"
)

func Test_node_Read(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&dfs.NodeConf{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	f := "/tmp/1.txt"
	utils.PanicIfErr(lfs.DeleteLocal(f))
	utils.PanicIfErr(lfs.CreateFileLocal(f, 32))
	utils.PanicIfErr(lfs.WriteLocal(f, 0, []byte("Hello world!")))

	bytes, err := n.Read("/1.txt", 0, int64(len("Hello world!")))
	utils.PanicIfErr(err)
	logger.Info(string(bytes))
	utils.PanicIfErr(lfs.DeleteLocal(f))
}
