package nodeimpl

import (
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
	"time"
)

func TestImpl_Stat(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(n.Delete("1.txt"))
	utils.PanicIfErr(n.CreateFile("1.txt", 10))
	stat, err := n.Stat("1.txt")
	utils.PanicIfErr(err)
	logger.Info(stat)

	utils.PanicIfErr(n.Delete("1.txt"))
}

func TestImpl_Stat2(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(n.Delete("1.txt"))
	stat, err := n.Stat("1.txt")
	utils.PanicIfErr(err)
	utils.PanicIf(stat.Exist())
}

func TestImpl_Exist(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(n.Delete("1.txt"))
	utils.PanicIfErr(n.CreateFile("1.txt", 10))
	stat, err := n.Stat("1.txt")
	utils.PanicIfErr(err)
	logger.Info(stat)
	exist, err := n.Exist("1.txt")
	utils.PanicIfErr(err)
	utils.PanicIf(!exist)

	utils.PanicIfErr(n.Delete("1.txt"))
}

func TestImpl_Exist2(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(n.Delete("1.txt"))
	stat, err := n.Stat("1.txt")
	logger.Info(err)
	utils.PanicIfErr(err)
	utils.PanicIf(stat.Exist())
	exist, err := n.Exist("1.txt")
	utils.PanicIfErr(err)
	utils.PanicIf(exist)
}
