package dfile

import (
	"github.com/madokast/GoDFS/internal/fs"
	"github.com/madokast/GoDFS/utils/serializer"
)

type Location struct {
	IP      string `json:"ip,omitempty"`
	Port    uint16 `json:"port,omitempty"`
	RootDir string `json:"rootDir,omitempty"`
}

// Meta 文件信息
type Meta interface {
	fs.Meta
	Locations() []*Location
}

// MetaImpl 文件信息，用于 Stat 请求 response
type MetaImpl struct {
	FullName_    string      `json:"fullName,omitempty"`
	Exist_       bool        `json:"exist,omitempty"`
	Size_        int64       `json:"size,omitempty"`
	IsDirectory_ bool        `json:"isDirectory,omitempty"`
	ModifyTime_  int64       `json:"modifyTime,omitempty"`
	Locations_   []*Location `json:"locations,omitempty"`
}

func (m *MetaImpl) FullName() string {
	return m.FullName_
}

func (m *MetaImpl) Exist() bool {
	return m.Exist_
}

func (m *MetaImpl) Size() int64 {
	return m.Size_
}

func (m *MetaImpl) IsDirectory() bool {
	return m.IsDirectory_
}

func (m *MetaImpl) Locations() []*Location {
	return m.Locations_
}

func (m *MetaImpl) ModifyTime() int64 {
	return m.ModifyTime_
}

func (m *MetaImpl) String() string {
	return serializer.JsonString(m)
}

func (l *Location) String() string {
	return serializer.JsonString(l)
}
