// Package lfs 本地文件操作
package lfs

import (
	"errors"
	"fmt"
	"github.com/madokast/GoDFS/utils/logger"
	"io"
	"os"
	path2 "path"
)

func ReadLocal(path string, offset, length int64) ([]byte, error) {
	f, err := os.OpenFile(path, os.O_RDONLY, 0666)
	if err != nil {
		return nil, errors.New("Open local " + path + " because " + err.Error())
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.Error("Close local file", err)
		}
	}(f)
	_, err = f.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, errors.New("Seek local " + path + " because " + err.Error())
	}
	data := make([]byte, length)
	n, err := f.Read(data)
	if err != nil {
		return nil, errors.New("Read local " + path + " because " + err.Error())
	}
	return data[:n], nil
}

func WriteLocal(path string, offset int64, data []byte) error {
	f, err := os.OpenFile(path, os.O_WRONLY, 0666)
	if err != nil {
		return errors.New("Open local " + path + " because " + err.Error())
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.Error("Close local file", err)
		}
	}(f)
	_, err = f.Seek(offset, io.SeekStart)
	if err != nil {
		return errors.New("Seek local " + path + " because " + err.Error())
	}
	n, err := f.Write(data)
	if err != nil {
		return errors.New("Write local " + path + " because " + err.Error())
	}
	if n != len(data) {
		return errors.New("Write local " + path + fmt.Sprintf(" %d bytes but success %d bytes", len(data), n))
	}
	return nil
}

// CreateFileLocal 创建文件，如果文件夹不存在，则创建
func CreateFileLocal(path string, size int64) error {
	father := path2.Dir(path)
	err := MkdirAllLocal(father)
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return errors.New("Create local " + path + " because " + err.Error())
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			logger.Error("Close local file", err)
		}
	}(file)
	_, err = file.WriteAt([]byte{0}, size-1)
	if err != nil {
		return errors.New("Create local " + path + " because " + err.Error())
	}
	return nil
}

func MkdirAllLocal(path string) error {
	return os.MkdirAll(path, 0666)
}

func DeleteLocal(path string) error {
	return os.RemoveAll(path)
}

// ListFilesLocal 列出本地目录下的所有文件/子目录
func ListFilesLocal(dir string) (files []string, dirs []string, err error) {
	paths, err := os.ReadDir(dir)
	if err != nil {
		return nil, nil, err
	}

	for _, path := range paths {
		if path.IsDir() {
			dirs = append(dirs, path2.Join(dir, path.Name()))
		} else {
			files = append(dirs, path2.Join(dir, path.Name()))
		}
	}
	return files, dirs, nil
}

func StatLocal(path string) (os.FileInfo, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("File " + path + " does not exist")
		}
		return nil, err
	}
	return stat, nil
}

func ExistLocal(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		logger.Error(err)
		return false
	}
	return true
}
