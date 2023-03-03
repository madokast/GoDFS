package locallock

import (
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
)

func TestNew(t *testing.T) {
	d1 := New()
	d2 := New()
	logger.Info(d1)
	logger.Info(d2)
	utils.PanicIf(d1 != d2, d1, d2)
}
