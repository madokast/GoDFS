package fsimpl

import (
	"fmt"
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

func TestImpl_ListFiles(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))

	DFS := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 10, FileReplicaNum: 1})
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

		for i := 0; i < 10; i++ {
			ex, err := DFS.Exist(fmt.Sprintf("/%00d.txt", i))
			utils.PanicIfErr(err)
			utils.PanicIf(!ex)
		}

		files, dirs, err := DFS.ListFiles("/")
		utils.PanicIfErr(err)
		logger.Info(files)
		logger.Info(dirs)
		utils.PanicIf(len(files) != 10)
		utils.PanicIf(len(dirs) != 0)
	}

	utils.PanicIfErr(DFS.Delete("/"))
}

func TestImpl_ListFiles1(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))

	DFS := New(locallock.New(), &fs.Conf{HashCircleReplicaNum: 10, FileReplicaNum: 2})
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

		for i := 0; i < 10; i++ {
			ex, err := DFS.Exist(fmt.Sprintf("/%00d.txt", i))
			utils.PanicIfErr(err)
			utils.PanicIf(!ex)
		}

		files, dirs, err := DFS.ListFiles("/")
		utils.PanicIfErr(err)
		logger.Info(files)
		logger.Info(dirs)
		utils.PanicIf(len(files) != 10)
		utils.PanicIf(len(dirs) != 0)
	}

	utils.PanicIfErr(DFS.Delete("/"))
}
