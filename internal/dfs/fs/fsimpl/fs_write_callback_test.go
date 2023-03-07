package fsimpl

import (
	"github.com/madokast/GoDFS/internal/dfs/fs"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/dfs/node/nodeimpl"
	"github.com/madokast/GoDFS/internal/dlock/locallock"
	"github.com/madokast/GoDFS/internal/fs/lfs"
	fs2 "github.com/madokast/GoDFS/internal/fs/write_callback"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"strconv"
	"testing"
	"time"
)

func TestImpl_RegisterWriteCallback(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))

	DFS := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 10, FileReplicaNum: 3})
	var nodes []node.Node
	for i := 0; i < 5; i++ {
		n := nodeimpl.New(&node.Info{
			IP:      "127.0.0.1",
			Port:    httputils.GetFreePort(),
			RootDir: "/tmp/dfs" + strconv.Itoa(i),
		})
		n.ListenAndServeGo()
		nodes = append(nodes, n)
		DFS.AddNode(n)
	}
	time.Sleep(100 * time.Millisecond)
	DFS.RefreshAliveNodesAndHashCircle()

	{
		path := "/table/col01.dat"
		utils.PanicIfErr(DFS.CreateFile(path, 16))

		DFS.RegisterWriteCallback(&fs2.Entry{
			FileName: path,
			Offset:   2,
			Length:   3,
			Callback: func() {
				logger.Info("do callback")
			},
		})

		utils.PanicIfErr(DFS.Write(path, 0, []byte("Hello, world!")))
	}

	utils.PanicIfErr(DFS.Delete("/"))
}

func TestImpl_RegisterWriteCallback2(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))

	DFS := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 10, FileReplicaNum: 3})
	DFS2 := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 10, FileReplicaNum: 3})
	var nodes []node.Node
	for i := 0; i < 5; i++ {
		n := nodeimpl.New(&node.Info{
			IP:      "127.0.0.1",
			Port:    httputils.GetFreePort(),
			RootDir: "/tmp/dfs" + strconv.Itoa(i),
		})
		n.ListenAndServeGo()
		nodes = append(nodes, n)
		DFS.AddNode(n)
	}
	DFS2.AddNode(nodes[1])
	for _, n := range nodes {
		DFS2.AddNode(n)
	}

	time.Sleep(100 * time.Millisecond)
	DFS.RefreshAliveNodesAndHashCircle()
	DFS2.RefreshAliveNodesAndHashCircle()

	{
		path := "/table/col01.dat"
		utils.PanicIfErr(DFS.CreateFile(path, 16))

		DFS.RegisterWriteCallback(&fs2.Entry{
			FileName: path,
			Offset:   2,
			Length:   3,
			Callback: func() {
				logger.Info("do callback 1")
			},
		})

		DFS2.RegisterWriteCallback(&fs2.Entry{
			FileName: path,
			Offset:   2,
			Length:   5,
			Callback: func() {
				logger.Info("do callback 2")
			},
		})

		utils.PanicIfErr(DFS.Write(path, 0, []byte("Hello, world!")))
	}

	utils.PanicIfErr(DFS.Delete("/"))
}
