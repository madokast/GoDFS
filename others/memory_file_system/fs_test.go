package memory_file_system

import (
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
)

func TestMemFileSystem_CreateFile(t *testing.T) {
	fs := NewMemFS()
	name := "/a/b/1.txt"
	file, err := fs.CreateFile(name)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("好好学习"))
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}

	file, err = fs.OpenFile(name)
	if err != nil {
		panic(err)
	}
	logger.Info(file.ReadString(256))
	err = file.Close()
	if err != nil {
		panic(err)
	}
}

func TestMemFileSystem_Stat(t *testing.T) {
	fs := NewMemFS()
	name := "/a/b/1.txt"
	file, err := fs.CreateFile(name)
	if err != nil {
		panic(err)
	}

	stat, err := fs.Stat(name)
	if err != nil {
		panic(err)
	}
	logger.Info(name, stat)

	_, err = file.Write([]byte("好好学习"))
	if err != nil {
		panic(err)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}

	file, err = fs.OpenFile(name)
	if err != nil {
		panic(err)
	}
	logger.Info(file.ReadString(256))
	err = file.Close()
	if err != nil {
		panic(err)
	}

	stat, err = fs.Stat(name)
	if err != nil {
		panic(err)
	}
	logger.Info(name, stat)

	err = fs.DeleteFile(name)
	if err != nil {
		panic(err)
	}

	_, err = fs.Stat(name)
	if err == nil {
		panic(err)
	}
	logger.Info(name, err)
}

func TestMemFileSystem_RenameFile(t *testing.T) {
	fs := NewMemFS()
	name := "/a/b/1.txt"
	name2 := "/a/b/2.txt"

	file, err := fs.CreateFile(name)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("好好学习"))
	if err != nil {
		panic(err)
	}
	err = file.Offset(0)
	if err != nil {
		panic(err)
	}
	logger.Info(file.ReadString(256))
	err = file.Close()
	if err != nil {
		panic(err)
	}
	logger.Info(fs.Stat(name))
	logger.Info(fs.Stat(name2))

	err = fs.RenameFile(name, name2)
	if err != nil {
		panic(err)
	}

	logger.Info(fs.Stat(name))
	logger.Info(fs.Stat(name2))

	file, err = fs.OpenFile(name2)
	if err != nil {
		panic(err)
	}
	logger.Info(file.ReadString(256))

	err = file.Close()
	if err != nil {
		panic(err)
	}
}
