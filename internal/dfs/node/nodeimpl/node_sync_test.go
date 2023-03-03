package nodeimpl

import (
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
	"time"
)

func TestImpl_Sync(t *testing.T) {
	n1 := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/tmp/a",
	})
	n1.ListenAndServeGo()

	n2 := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/tmp/b",
	})
	n2.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(n1.Delete("1.txt"))
	utils.PanicIfErr(n2.Delete("1.txt"))

	utils.PanicIfErr(n1.CreateFile("1.txt", 10))
	utils.PanicIfErr(n1.Write("1.txt", 0, []byte("hello")))

	utils.PanicIfErr(n2.Sync(n1, "1.txt"))
	n2Read, err := n2.Read("1.txt", 0, 5)
	utils.PanicIfErr(err)
	logger.Info(string(n2Read))
	utils.PanicIf(string(n2Read) != "hello")

	utils.PanicIfErr(n1.Delete("1.txt"))
	utils.PanicIfErr(n2.Delete("1.txt"))
}

func TestImpl_Sync2(t *testing.T) {
	n1 := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/tmp/a",
	})
	n1.ListenAndServeGo()

	n2 := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/tmp/b",
	})
	n2.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(n1.Delete("1.txt"))
	utils.PanicIfErr(n2.Delete("1.txt"))

	utils.PanicIfErr(n1.CreateFile("1.txt", 10))
	utils.PanicIfErr(n1.Write("1.txt", 0, []byte("hello")))

	utils.PanicIfErr(n2.CreateFile("1.txt", 10))
	utils.PanicIfErr(n2.Write("1.txt", 0, []byte("hello")))

	utils.PanicIfErr(n2.Sync(n1, "1.txt"))
	n2Read, err := n2.Read("1.txt", 0, 5)
	utils.PanicIfErr(err)
	logger.Info(string(n2Read))
	utils.PanicIf(string(n2Read) != "hello")

	utils.PanicIfErr(n1.Delete("1.txt"))
	utils.PanicIfErr(n2.Delete("1.txt"))
}

func TestImpl_Sync3(t *testing.T) {
	n1 := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/tmp/a",
	})
	n1.ListenAndServeGo()

	n2 := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/tmp/b",
	})
	n2.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	utils.PanicIfErr(n1.Delete("1.txt"))
	utils.PanicIfErr(n2.Delete("1.txt"))

	utils.PanicIfErr(n1.CreateFile("1.txt", 10))
	utils.PanicIfErr(n1.Write("1.txt", 0, []byte("hello")))

	utils.PanicIfErr(n2.CreateFile("1.txt", 10))
	utils.PanicIfErr(n2.Write("1.txt", 0, []byte("sync")))

	utils.PanicIfErr(n2.Sync(n1, "1.txt"))
	n2Read, err := n2.Read("1.txt", 0, 5)
	utils.PanicIfErr(err)
	logger.Info(string(n2Read))
	utils.PanicIf(string(n2Read) != "hello")

	utils.PanicIfErr(n1.Delete("1.txt"))
	utils.PanicIfErr(n2.Delete("1.txt"))
}
