package main

import (
	"github.com/madokast/GoDFS/internal/dfs"
	"github.com/madokast/GoDFS/internal/dfs/fs"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/dlock/locallock"
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
)

func main() {
	node1 := dfs.NewNode(&node.Info{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/tmp/dfs1",
	})
	node1.ListenAndServeGo()

	node2 := dfs.NewNode(&node.Info{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/tmp/dfs2",
	})
	node2.ListenAndServeGo()

	DFS := dfs.NewDFS(locallock.New(), &fs.Conf{
		HashCircleReplicaNum: 3,
	})

	DFS.AddNode(node1)
	DFS.AddNode(node2)
	DFS.RefreshAliveNodesAndHashCircle()
	DFS.RefreshAliveNodesAndHashCircle()

	logger.Info(DFS)
	utils.PanicIf(len(DFS.AllNodes()) != 2)

}
