package lfs

import (
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
)

func TestListFilesLocal(t *testing.T) {
	logger.Info(ListFilesLocal("/a/a/b"))
}

func TestStatLocal(t *testing.T) {
	logger.Info(StatLocal("/a/b/v"))
}
