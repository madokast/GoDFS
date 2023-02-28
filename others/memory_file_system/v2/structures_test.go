package v2

import (
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
)

func Test_canReadMode(t *testing.T) {
	mode := writeMode
	logger.Info(canReadMode(mode))
}
