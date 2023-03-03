package fsimpl

import (
	"github.com/madokast/GoDFS/internal/dfs/fs"
	"github.com/madokast/GoDFS/internal/dfs/node"
	"github.com/madokast/GoDFS/internal/dfs/node/nodeimpl"
	"github.com/madokast/GoDFS/internal/dlock/locallock"
	"github.com/madokast/GoDFS/utils/httputils"
	"testing"
	"time"
)

func TestImpl_RefreshAliveNodesAndHandCircle(t *testing.T) {
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
		HashCircleReplicaNum: 3,
	})

	DFS.AddNode(node1)
	DFS.AddNode(node2)
	DFS.RefreshAliveNodesAndHashCircle()
}
