package allocator_impl

import (
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/logger"
	"strconv"
	"testing"
)

func TestCreateHeader(t *testing.T) {
	header := CreateHeader(63)
	logger.Info("h", strconv.FormatInt(int64(header.freeFlagSize), 2))
	logger.Info(&header)
}

func TestMemHeader_Size(t *testing.T) {
	header := CreateHeader(63)
	utils.PanicIf(header.Size() != 63)
}

func TestMemHeader_IsFree(t *testing.T) {
	header := CreateHeader(63)
	utils.PanicIf(header.IsFree())
}

func TestMemHeader_Free(t *testing.T) {
	header := CreateHeader(63)
	header.Free()
	utils.PanicIf(!header.IsFree())
	logger.Info(&header)
}

func TestMemHeader_Free2(t *testing.T) {
	header := CreateHeader(63)
	header.Free()
	logger.Info("== ERROR OK ==")
	header.Free()
}

func TestMemHeader_WriteHeaderTo(t *testing.T) {
	bytes := make([]byte, MemHeaderSz)
	header := CreateHeader(63)
	header.WriteHeaderTo(bytes)
	logger.Info(bytes)
}

func TestReadHeader(t *testing.T) {
	bytes := make([]byte, MemHeaderSz)
	header := CreateHeader(63)
	header.WriteHeaderTo(bytes)

	header2 := ReadHeader(bytes)
	utils.PanicIf(header != header2, header, header2)
}
