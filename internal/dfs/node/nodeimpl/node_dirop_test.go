package nodeimpl

import (
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
	"time"
)

func TestImpl_ListFiles(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(n.Delete("/abc"))
	utils.PanicIfErr(n.CreateFile("/abc/1.txt", 10))
	utils.PanicIfErr(n.CreateFile("/abc/2.txt", 10))
	utils.PanicIfErr(n.CreateFile("/abc/3.txt", 10))
	utils.PanicIfErr(n.CreateFile("/abc/4.txt", 10))
	utils.PanicIfErr(n.CreateFile("/abc/def1/21.txt", 10))
	utils.PanicIfErr(n.CreateFile("/abc/def2/22.txt", 10))
	utils.PanicIfErr(n.CreateFile("/abc/def3/23.txt", 10))
	utils.PanicIfErr(n.CreateFile("/abc/def4/24.txt", 10))

	files, dirs, err := n.ListFiles("/abc")
	utils.PanicIfErr(err)
	logger.Info(files)
	logger.Info(dirs)

	files, dirs, err = n.ListFiles("/abc/def2")
	utils.PanicIfErr(err)
	logger.Info(files)
	logger.Info(dirs)

	utils.PanicIfErr(n.Delete("/abc"))
}
