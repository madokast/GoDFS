package nodeimpl

import (
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/fs/lfs"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
	"time"
)

func Test_node_Read(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
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

func Test_node_Read2(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	f := "/tmp/1.txt"
	utils.PanicIfErr(lfs.DeleteLocal(f))

	_, err := n.Read("/1.txt", 0, int64(len("Hello world!")))
	utils.PanicIf(err == nil, err)
	utils.PanicIfErr(lfs.DeleteLocal(f))
}

func Test_node_Write(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	f := "/tmp/1.txt"
	utils.PanicIfErr(lfs.DeleteLocal(f))
	utils.PanicIfErr(lfs.CreateFileLocal(f, 32))

	utils.PanicIfErr(n.Write("/1.txt", 0, []byte("Hello, world!")))
	read, err := n.Read("/1.txt", 0, int64(len("Hello world!")))
	utils.PanicIfErr(err)
	logger.Info(string(read))
	utils.PanicIfErr(lfs.DeleteLocal(f))
}

func Test_node_Write2(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	f := "/tmp/1.txt"
	utils.PanicIfErr(lfs.DeleteLocal(f))
	utils.PanicIfErr(lfs.CreateFileLocal(f, 32))

	utils.PanicIfErr(n.Write("/1.txt", 0, []byte("Hello")))
	utils.PanicIfErr(n.Write("/1.txt", int64(len("Hello")), []byte(", world!")))
	read, err := n.Read("/1.txt", 0, int64(len("Hello world!")))
	utils.PanicIfErr(err)
	logger.Info(string(read))
	utils.PanicIfErr(lfs.DeleteLocal(f))
}

func Test_node_Write3(t *testing.T) {
	port := httputils.GetFreePort()
	n := New(&node.Info{
		IP:      "127.0.0.1",
		Port:    port,
		RootDir: "/tmp",
	})
	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	f := "/tmp/1.txt"
	utils.PanicIfErr(lfs.DeleteLocal(f))
	utils.PanicIfErr(lfs.CreateFileLocal(f, 32))

	utils.PanicIfErr(n.Write("/1.txt", 0, []byte("Hello")))
	utils.PanicIfErr(n.Write("/1.txt", int64(len("Hello")), []byte(", world!")))
	read1, err := n.Read("/1.txt", 0, int64(len("Hello world!"))/2)
	utils.PanicIfErr(err)
	read2, err := n.Read("/1.txt", int64(len("Hello world!"))/2, int64(len("Hello world!"))-int64(len("Hello world!"))/2)
	utils.PanicIfErr(err)
	logger.Info(string(read1))
	logger.Info(string(read2))
	utils.PanicIfErr(lfs.DeleteLocal(f))
}
