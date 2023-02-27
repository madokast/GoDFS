package v2

import (
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
)

func Test_split(t *testing.T) {
	// [/]
	logger.Info(split("/"))
	// [/ a]
	logger.Info(split("/a"))
	// [/ a b]
	logger.Info(split("/a/b"))
	// [/ a b 1.txt]
	logger.Info(split("/a/b/1.txt"))
}

func Test_hierarchyFullPath(t *testing.T) {
	// [/]
	logger.Info(hierarchyFullPath("/"))
	// [/ a]
	logger.Info(hierarchyFullPath("/a"))
	// [/ a b]
	logger.Info(hierarchyFullPath("/a/b"))
	// [/ a b 1.txt]
	logger.Info(hierarchyFullPath("/a/b/1.txt"))
}

func TestMemFileSystem2_ListAllPath(t *testing.T) {
	fs := NewMemFS()
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_Stat(t *testing.T) {
	fs := NewMemFS()
	logger.Info(fs.Stat("/"))
}

func TestMemFileSystem2_MakeDirectory(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectory("/a")
	if err != nil {
		panic(err)
	}
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_MakeDirectory2(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectory("/a")
	if err != nil {
		panic(err)
	}
	err = fs.MakeDirectory("/b")
	if err != nil {
		panic(err)
	}
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_MakeDirectory3(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectory("/a")
	if err != nil {
		panic(err)
	}
	err = fs.MakeDirectory("/a/b")
	if err != nil {
		panic(err)
	}
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}
