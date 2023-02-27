package v2

import (
	"errors"
	"fmt"
)

/*============================= FileIO ==========================================*/

func (m *MemFile2) Read(p []byte) (n int, err error) {
	n = copy(p, m.meta.data[m.local:])
	m.local += int64(n)
	return n, nil
}

func (m *MemFile2) Write(p []byte) (n int, err error) {
	n = copy(m.meta.data[m.local:], p)
	m.local += int64(n)
	return n, nil
}

func (m *MemFile2) Offset(offset int64) error {
	if offset < 0 || offset >= m.meta.size {
		return errors.New(fmt.Sprintf("offset %d invaid for file %s size %d", offset, m.meta.fullName, m.meta.size))
	}
	m.local = offset
	return nil
}

// ReadString impl from Read
func (m *MemFile2) ReadString(limit int) (string, error) {
	data := make([]byte, limit)
	n, err := m.Read(data)
	if err != nil {
		return "", err
	}
	return string(data[:n]), nil
}

func (m *MemFile2) Close() error {
	return nil
}

/*============================= File Meta ==========================================*/

func (m *MemFileMeta2) FullName() string {
	return m.fullName
}

func (m *MemFileMeta2) Size() int64 {
	return m.size
}

func (m *MemFileMeta2) IsDirectory() bool {
	return m.isDirectory
}
