package memory_file_system

import (
	"errors"
	"github.com/madokast/GoDFS/internal/ifile"
	"strings"
)

func (m *MemFileSystem) CreateFile(name string) (ifile.FileIO, error) {
	_, ok := m.fileData[name]
	if ok {
		return nil, errors.New("File " + name + " already exists")
	}
	bytes := make([]byte, 0)
	m.fileData[name] = bytes
	return &MemFile{data: bytes, fullName: name, local: 0, fs: m}, nil
}

func (m *MemFileSystem) OpenFile(name string) (ifile.FileIO, error) {
	data, ok := m.fileData[name]
	if !ok {
		return nil, errors.New("File " + name + " does not exist")
	}
	return &MemFile{data: data, fullName: name, local: 0, fs: m}, nil
}

func (m *MemFileSystem) DeleteFile(name string) error {
	_, ok := m.fileData[name]
	if !ok {
		return errors.New("File " + name + " does not exist")
	}
	delete(m.fileData, name)
	return nil
}

func (m *MemFileSystem) RenameFile(name, newName string) error {
	data, ok := m.fileData[name]
	if !ok {
		return errors.New("File " + name + " does not exist")
	}

	_, ok = m.fileData[newName]
	if ok {
		return errors.New("File " + newName + " already exists")
	}

	delete(m.fileData, name)
	m.fileData[newName] = data
	return nil
}

func (m *MemFileSystem) MakeDirectory(string) error {
	panic("implement me")
}

func (m *MemFileSystem) MakeDirectories(string) error {
	panic("implement me")
}

func (m *MemFileSystem) DeleteDirectory(string) error {
	panic("implement me")
}

func (m *MemFileSystem) DeleteDirectories(string) error {
	panic("implement me")
}

func (m *MemFileSystem) ReadDirectory(dir string) ([]ifile.FileMeta, error) {
	stat, err := m.Stat(dir)
	if err != nil {
		return nil, err
	}
	if !stat.IsDirectory() {
		return nil, errors.New(dir + " is not a dir")
	}

	dirLen := len(dir)
	files := make([]ifile.FileMeta, 0)
	for file := range m.fileData {
		if len(file) > dirLen+1 && file[dirLen] == '/' && strings.HasPrefix(file, dir) {
			sub := file[dirLen:]
			if strings.Count(sub, "/") == 1 {
				stat, err := m.Stat(file)
				if err != nil {
					panic(err)
				}
				files = append(files, stat)
			}
		}
	}
	return files, nil
}

func (m *MemFileSystem) RenameDirectory(string, string) error {
	panic("implement me")
}

func (m *MemFileSystem) Stat(name string) (ifile.FileMeta, error) {
	data, ok := m.fileData[name]
	split := strings.Split(name, "/")
	if ok {
		return &MemFileMeta{
			baseName:    split[len(split)-1],
			fullName:    name,
			size:        int64(len(data)),
			isDirectory: false,
		}, nil
	}

	dirLen := len(name)
	for file := range m.fileData {
		if len(file) > dirLen+1 && file[dirLen] == '/' && strings.HasPrefix(file, name) {
			return &MemFileMeta{
				baseName:    split[len(split)-1],
				fullName:    name,
				size:        int64(len(data)),
				isDirectory: true,
			}, nil
		}
	}
	return nil, errors.New(name + " does not exist")
}
