package nodeimpl

import (
	"github.com/madokast/GoDFS/internal/dfs/lfs"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
	"time"
)

func Test_node_CreateFile(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(lfs.DeleteLocal("/tmp/1.txt"))
	utils.PanicIfErr(n.CreateFile("/1.txt", 16))
	exist := lfs.ExistLocal("/tmp/1.txt")
	utils.PanicIf(!exist, "???")

	utils.PanicIfErr(lfs.DeleteLocal("/tmp/1.txt"))
}

func Test_node_CreateFile2(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(lfs.DeleteLocal("/tmp/1.txt"))
	utils.PanicIfErr(n.CreateFile("/1.txt", 16))
	exist := lfs.ExistLocal("/tmp/1.txt")
	utils.PanicIf(!exist, "???")

	utils.PanicIfErr(n.Write("/1.txt", 3, []byte("abc")))
	read, err := n.Read("/1.txt", 3, 3)
	utils.PanicIfErr(err)
	utils.PanicIf(string(read) != "abc", string(read))
	logger.Info(string(read))

	utils.PanicIfErr(lfs.DeleteLocal("/tmp/1.txt"))
}

func Test_node_Delete(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(n.Delete("/1.txt"))
	utils.PanicIfErr(n.Delete("/1.txt"))

	ex := lfs.ExistLocal("/tmp/1.txt")
	utils.PanicIf(ex)
}

func Test_node_Delete2(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(n.Delete("/1.txt"))
	utils.PanicIfErr(n.CreateFile("1.txt", 10))
	ex := lfs.ExistLocal("/tmp/1.txt")
	utils.PanicIf(!ex)

	utils.PanicIfErr(n.Delete("/1.txt"))
	ex = lfs.ExistLocal("/tmp/1.txt")
	utils.PanicIf(ex)
}

func Test_node_Delete3(t *testing.T) {
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
	utils.PanicIfErr(n.Write("1.txt", 2, []byte("abc")))
	read, err := n.Read("1.txt", 2, 3)
	utils.PanicIfErr(err)
	utils.PanicIf(string(read) != "abc")
	logger.Info(string(read))

	utils.PanicIfErr(n.Delete("1.txt"))
	utils.PanicIf(lfs.ExistLocal("/tmp/1.txt"))
}

func TestImpl_MD5(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(n.Delete("1.txt"))
	utils.PanicIfErr(n.Delete("2.txt"))
	utils.PanicIfErr(n.CreateFile("1.txt", 10))
	utils.PanicIfErr(n.CreateFile("2.txt", 10))
	utils.PanicIfErr(n.Write("1.txt", 2, []byte("abc")))
	utils.PanicIfErr(n.Write("2.txt", 2, []byte("abc")))

	md51, err := n.MD5("1.txt")
	utils.PanicIfErr(err)
	md52, err := n.MD5("2.txt")
	utils.PanicIfErr(err)
	logger.Info(md51)
	logger.Info(md52)
	utils.PanicIf(md51 != md52)

	utils.PanicIfErr(n.Delete("1.txt"))
	utils.PanicIfErr(n.Delete("2.txt"))
}
