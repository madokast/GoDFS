package fsimpl

import (
	"fmt"
	"github.com/madokast/GoDFS/internal/dfs/fs"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/dfs/node/nodeimpl"
	"github.com/madokast/GoDFS/internal/dlock/locallock"
	"github.com/madokast/GoDFS/internal/fs/lfs"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"strconv"
	"testing"
	"time"
)

func TestImpl_RefreshAliveNodesAndHashCircle(t *testing.T) {
	node1 := nodeimpl.New(&node.Info{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/tmp/dfs1",
	})
	node1.ListenAndServeGo()

	node2 := nodeimpl.New(&node.Info{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/tmp/dfs2",
	})
	node2.ListenAndServeGo()

	time.Sleep(100 * time.Millisecond)

	DFS := New(locallock.New(), &fs.Conf{
		HashCircleReplicaNum: 3, FileReplicaNum: 1,
	})

	DFS.AddNode(node1)
	DFS.AddNode(node2)
	DFS.RefreshAliveNodesAndHashCircle()
}

func TestNew(t *testing.T) {
	node1 := nodeimpl.New(&node.Info{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/tmp/dfs1",
	})
	node1.ListenAndServeGo()

	time.Sleep(100 * time.Millisecond)

	DFS := New(locallock.New(), &fs.Conf{
		HashCircleReplicaNum: 3, FileReplicaNum: 1,
	})

	DFS.AddNode(node1)
	DFS.RefreshAliveNodesAndHashCircle()

	nodes := DFS.HashDistribute("/1.txt", 1)
	logger.Info(nodes)

	nodes = DFS.HashDistribute("/1.txt", 2)
	logger.Info(nodes)
}

func TestImpl_RefreshAliveNodesAndHashCircle1(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))

	DFS := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 1, FileReplicaNum: 4})
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
	for _, n := range nodes {
		utils.PanicIfErr(n.CreateFile("/1.txt", 10))
	}

	{
		newNode := nodeimpl.New(&node.Info{
			IP:      "127.0.0.1",
			Port:    httputils.GetFreePort(),
			RootDir: "/tmp/dfs" + strconv.Itoa(4),
		})
		newNode.ListenAndServeGo()

		time.Sleep(100 * time.Millisecond)

		DFS.AddNode(newNode)
		DFS.RefreshAliveNodesAndHashCircle()

		existFile, err := newNode.Exist("/1.txt")
		utils.PanicIfErr(err)
		utils.PanicIf(!existFile)
	}

	utils.PanicIfErr(DFS.Delete("/"))
}

func TestImpl_RefreshAliveNodesAndHashCircle2(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))

	DFS := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 1, FileReplicaNum: 4})
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
	for _, n := range nodes {
		utils.PanicIfErr(n.CreateFile("/1.txt", 10))
		utils.PanicIfErr(n.Write("/1.txt", 2, []byte("xyn")))
	}

	{
		newNode := nodeimpl.New(&node.Info{
			IP:      "127.0.0.1",
			Port:    httputils.GetFreePort(),
			RootDir: "/tmp/dfs" + strconv.Itoa(4),
		})
		newNode.ListenAndServeGo()

		time.Sleep(100 * time.Millisecond)

		DFS.AddNode(newNode)
		DFS.RefreshAliveNodesAndHashCircle()

		read, err := newNode.Read("/1.txt", 2, 3)
		utils.PanicIfErr(err)
		utils.PanicIf(string(read) != "xyn")
		logger.Info(string(read))
	}

	utils.PanicIfErr(DFS.Delete("/"))
}

func TestImpl_RefreshAliveNodesAndHashCircle3(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))

	DFS := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 1, FileReplicaNum: 3})
	var nodes []node.Node
	for i := 0; i < 4; i++ {
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
		utils.PanicIfErr(DFS.CreateFile("/1.txt", 10))
		closeOne := true
		for _, n := range nodes {
			ex, err := n.Exist("/1.txt")
			utils.PanicIfErr(err)
			if ex {
				logger.Info("exist", n.Key())
				if closeOne {
					n.Close()
					closeOne = false
				}
			} else {
				logger.Info("un-exist", n.Key())
			}
		}

		// 应该刷新到 un-exist node 中
		logger.Info("=== WARN OK ===")
		DFS.RefreshAliveNodesAndHashCircle()

		ex, err := DFS.Exist("/1.txt")
		utils.PanicIfErr(err)
		utils.PanicIf(!ex)
	}

	utils.PanicIfErr(DFS.Delete("/"))
}

func TestImpl_RefreshAliveNodesAndHashCircle4(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))

	DFS := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 5, FileReplicaNum: 3})
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

	{
		for i := 0; i < 10; i++ {
			utils.PanicIfErr(DFS.CreateFile(fmt.Sprintf("/%00d.txt", i), 10))
		}
		newNode := nodeimpl.New(&node.Info{
			IP:      "127.0.0.1",
			Port:    httputils.GetFreePort(),
			RootDir: "/tmp/dfs" + strconv.Itoa(5),
		})
		newNode.ListenAndServeGo()
		time.Sleep(100 * time.Millisecond)

		DFS.AddNode(newNode)
		DFS.RefreshAliveNodesAndHashCircle()

		for i := 0; i < 10; i++ {
			ex, err := DFS.Exist(fmt.Sprintf("/%00d.txt", i))
			utils.PanicIfErr(err)
			utils.PanicIf(!ex)
			logger.Info(fmt.Sprintf("/%00d.txt", i), "ex", ex)
		}
	}

	utils.PanicIfErr(DFS.Delete("/"))
}
