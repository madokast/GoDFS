package node

import (
	"github.com/madokast/GoDFS/internal/dfs"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	n := New(&dfs.NodeConf{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/dfs",
	})

	logger.Info("=== error ok ===")
	alive := n.Ping()
	if alive {
		panic(alive)
	}
	logger.Info(alive)
}

func TestAlive(t *testing.T) {
	n := New(&dfs.NodeConf{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/dfs",
	})

	n.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	alive := n.Ping()
	if !alive {
		panic(alive)
	}
	logger.Info(alive)
}

func TestAlive2(t *testing.T) {
	n1 := New(&dfs.NodeConf{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/dfs1",
	})

	n1.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	n2 := New(&dfs.NodeConf{
		IP:      "127.0.0.1",
		Port:    httputils.GetFreePort(),
		RootDir: "/dfs2",
	})

	n2.ListenAndServeGo()
	time.Sleep(100 * time.Millisecond)

	alive := n1.Ping()
	if !alive {
		panic(alive)
	}
	logger.Info(alive)

	alive = n2.Ping()
	if !alive {
		panic(alive)
	}
	logger.Info(alive)
}
