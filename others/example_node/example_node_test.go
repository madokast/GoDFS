package example_node

import (
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/jsonutils"
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
	"time"
)

func Test_node_getInfo(t *testing.T) {
	n := node{
		IP:   "127.0.0.1",
		Port: httputils.GetFreePort(),
		Info: "example_node",
	}

	n.serverGo()
	time.Sleep(100 * time.Millisecond)

	info, err := n.getInfo()
	if err != nil {
		logger.Error(err)
	}
	logger.Info(jsonutils.String(info))
}
