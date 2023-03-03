package fsimpl

import (
	"github.com/madokast/GoDFS/internal/dfs/fs"
	"github.com/madokast/GoDFS/internal/dfs/lfs"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/dfs/node/nodeimpl"
	"github.com/madokast/GoDFS/internal/dlock/locallock"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"strconv"
	"testing"
	"time"
)

func TestImpl_ExistFile(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))

	DFS := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 1, FileReplicaNum: 3})
	var nodes []node.Node
	for i := 0; i < 3; i++ {
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

	existFile, err := DFS.Exist("/1.txt")
	utils.PanicIfErr(err)
	utils.PanicIf(existFile)
	utils.PanicIfErr(DFS.Delete("/"))
}

func TestImpl_ExistFile2(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))

	DFS := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 1, FileReplicaNum: 3})
	var nodes []node.Node
	for i := 0; i < 3; i++ {
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

	utils.PanicIfErr(nodes[0].CreateFile("/1.txt", 10))

	logger.Info("=== WARN OK ===")
	existFile, err := DFS.Exist("/1.txt")
	utils.PanicIfErr(err)
	utils.PanicIf(existFile)
	utils.PanicIfErr(DFS.Delete("/"))
}

func TestImpl_ExistFile3(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))
	DFS := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 1, FileReplicaNum: 3})
	var nodes []node.Node
	for i := 0; i < 3; i++ {
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

	utils.PanicIfErr(nodes[0].CreateFile("/1.txt", 10))
	utils.PanicIfErr(nodes[1].CreateFile("/1.txt", 10))

	logger.Info("=== WARN OK ===")
	existFile, err := DFS.Exist("/1.txt")
	utils.PanicIfErr(err)
	utils.PanicIf(!existFile)
	utils.PanicIfErr(DFS.Delete("/"))
}
