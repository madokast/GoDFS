package allocator_impl

import (
	"github.com/madokast/GoDFS/internal/dfs"
	"github.com/madokast/GoDFS/internal/dfs/fs"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/dlock/locallock"
	"github.com/madokast/GoDFS/internal/fs/lfs"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"strconv"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	utils.PanicIfErr(lfs.DeleteLocal("/tmp"))
	DFS := dfs.NewDFS(locallock.New(), &fs.Conf{HashCircleReplicaNum: 10, FileReplicaNum: 3})
	var nodes []node.Node
	for i := 0; i < 5; i++ {
		n := dfs.NewNode(&node.Info{
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
		logger.Info("Starting test.")
		allocator, err := New(DFS, DFS, "/mm", 256, 128)
		utils.PanicIfErr(err)
		pointer, err := allocator.Allocate(5, false)
		utils.PanicIfErr(err)
		err = allocator.Write(pointer, 0, []byte("Hello"))
		utils.PanicIfErr(err)
		logger.Info(allocator.ReadBytes(pointer))
		readString, err := allocator.ReadString(pointer)
		utils.PanicIfErr(err)
		logger.Info(readString)
		utils.PanicIf(readString != "Hello")
	}

	utils.PanicIfErr(DFS.Delete("/"))
}
