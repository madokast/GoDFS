package v2

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/madokast/GoDFS/internal/ifile"
	"github.com/madokast/GoDFS/utils/logger"
	"sort"
	"strings"
)

/**
所有路径名必须是全称 /a/b/c 形式
*/

func (fs *MemFileSystem2) CreateFile(name string, size int64) (ifile.FileIO, error) {
	//TODO implement me
	panic("implement me")
}

func (fs *MemFileSystem2) OpenFile(name string) (ifile.FileIO, error) {
	meta, err := fs.Stat(name)
	if err != nil {
		return nil, err
	}

	if meta.IsDirectory() {
		return nil, errors.New(name + " is a dir.")
	}

	return &MemFile2{
		local: 0,
		fs:    fs,
		meta:  meta.(*MemFileMeta2),
	}, nil
}

func (fs *MemFileSystem2) DeleteFile(name string) error {
	//TODO implement me
	panic("implement me")
}

func (fs *MemFileSystem2) RenameFile(name, newName string) error {
	//TODO implement me
	panic("implement me")
}

func (fs *MemFileSystem2) MakeDirectory(name string) error {
	hPaths, err := hierarchyFullPath(name)
	if err != nil {
		return err
	}
	if len(hPaths) == 1 {
		logger.Warn(name, "already exists when making directory.")
		return nil
	}
	fatherPath := hPaths[len(hPaths)-2]
	father, err := fs.Stat(fatherPath)
	if err != nil {
		return err
	}
	if !father.IsDirectory() {
		return errors.New("Cannot mkdir " + name + " because " + fatherPath + " is not a dir.")
	}
	fileMeta2 := father.(*MemFileMeta2)

	dirName := hPaths[len(hPaths)-1]
	fileMeta2.objects = append(fileMeta2.objects, dirName)
	sort.Strings(fileMeta2.objects)
	fs.fileMap[dirName] = &MemFileMeta2{
		fullName:    dirName,
		size:        0,
		isDirectory: true,
		data:        nil,
		objects:     []string{},
	}
	return nil
}

func (fs *MemFileSystem2) MakeDirectories(name string) error {
	//TODO implement me
	panic("implement me")
}

func (fs *MemFileSystem2) DeleteDirectory(name string) error {
	//TODO implement me
	panic("implement me")
}

func (fs *MemFileSystem2) DeleteDirectoryAll(name string) error {
	//TODO implement me
	panic("implement me")
}

func (fs *MemFileSystem2) ReadDirectory(dir string) ([]ifile.FileMeta, error) {
	//TODO implement me
	panic("implement me")
}

func (fs *MemFileSystem2) RenameDirectory(name, newName string) error {
	//TODO implement me
	panic("implement me")
}

func (fs *MemFileSystem2) Stat(name string) (ifile.FileMeta, error) {
	hPaths, err := hierarchyFullPath(name)
	if err != nil {
		return nil, err
	}
	var meta *MemFileMeta2
	for _, hPath := range hPaths {
		m, ok := fs.fileMap[hPath]
		if ok {
			meta = m
		} else {
			return nil, errors.New(fmt.Sprintf("%s does not exist because /%s does not exist.", name, hPath))
		}
	}
	if meta == nil {
		panic("??")
	}
	return meta, nil
}

func (fs *MemFileSystem2) ListAllPath() []ifile.FileMeta {
	var ret []ifile.FileMeta
	stack := list.New()
	stack.PushFront("/")
	for stack.Len() > 0 {
		pop := stack.Front()
		stack.Remove(pop)
		popMeta := fs.fileMap[pop.Value.(string)]
		ret = append(ret, popMeta)
		if popMeta.IsDirectory() {
			for i := len(popMeta.objects) - 1; i >= 0; i-- {
				stack.PushFront(popMeta.objects[i])
			}
		}
	}
	return ret
}

// 路径切分
func split(name string) ([]string, error) {
	if !strings.HasPrefix(name, "/") {
		return nil, errors.New("Invalid path " + name)
	}
	if name == "/" {
		return []string{"/"}, nil
	}

	if strings.HasSuffix(name, "/") {
		return nil, errors.New("Invalid path " + name)
	}

	return append([]string{"/"}, strings.Split(name, "/")[1:]...), nil
}

func hierarchyFullPath(name string) ([]string, error) {
	bases, err := split(name)
	if err != nil {
		return nil, err
	}
	ret := []string{"/"}
	fullPath := ""
	for i := 1; i < len(bases); i++ {
		fullPath += "/" + bases[i]
		ret = append(ret, fullPath)
	}
	return ret, nil
}
