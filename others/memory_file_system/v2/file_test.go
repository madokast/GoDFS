package v2

import (
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
)

func TestMemFile2_Close(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/1.txt", 4)
	utils.PanicIfErr(err)
	fileIO, err := fs.OpenFile("/1.txt", true)
	utils.PanicIfErr(err)
	err = fs.DeleteFile("/1.txt")
	logger.Info(err)
	utils.PanicIf(err == nil, err)

	err = fileIO.Close()
	utils.PanicIfErr(err)
	logger.Info("=== ok warn ===")
	err = fileIO.Close() // warn
	logger.Info("=== ok warn ===")
	utils.PanicIfErr(err)
	exist := fs.Exist("/1.txt")
	utils.PanicIf(!exist)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta)
	}

	err = fs.DeleteFile("/1.txt")
	utils.PanicIfErr(err)
}
