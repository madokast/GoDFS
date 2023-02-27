package v1

import (
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
)

func TestMemFile_Offset(t *testing.T) {
	fs := NewMemFS()
	name := "/a/b/1.txt"
	file, err := fs.CreateFile(name, 4)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("aaaa"))
	if err != nil {
		panic(err)
	}
	err = file.Offset(1)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("b"))
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
	readString, err := file.ReadString(256)
	if err != nil {
		panic(err)
	}
	logger.Info(readString)
	if readString != "abaa" {
		panic(readString)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
}

func TestMemFile_Offset2(t *testing.T) {
	fs := NewMemFS()
	name := "/a/b/1.txt"
	file, err := fs.CreateFile(name, 4)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("aaaa"))
	if err != nil {
		panic(err)
	}
	err = file.Offset(1)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("bd"))
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
	readString, err := file.ReadString(256)
	if err != nil {
		panic(err)
	}
	logger.Info(readString)
	if readString != "abda" {
		panic(readString)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
}

func TestMemFile_Offset3(t *testing.T) {
	fs := NewMemFS()
	name := "/a/b/1.txt"
	file, err := fs.CreateFile(name, 4)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("aaaa"))
	if err != nil {
		panic(err)
	}
	err = file.Offset(1)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("bbb"))
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
	readString, err := file.ReadString(256)
	if err != nil {
		panic(err)
	}
	logger.Info(readString)
	if readString != "abbb" {
		panic(readString)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
}

func TestMemFile_Offset4(t *testing.T) {
	fs := NewMemFS()
	name := "/a/b/1.txt"
	file, err := fs.CreateFile(name, 5)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("aaaa"))
	if err != nil {
		panic(err)
	}
	err = file.Offset(1)
	if err != nil {
		panic(err)
	}
	_, err = file.Write([]byte("bbbb"))
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
	readString, err := file.ReadString(256)
	if err != nil {
		panic(err)
	}
	logger.Info(readString)
	if readString != "abbbb" {
		panic(readString)
	}
	err = file.Close()
	if err != nil {
		panic(err)
	}
}
