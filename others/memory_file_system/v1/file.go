package v1

import (
	"errors"
	"fmt"
	"github.com/madokast/GoDFS/utils/logger"
	"strconv"
)

func (m *MemFile) Close() error {
	_, ok := m.fs.fileData[m.fullName]
	if !ok {
		return errors.New("Close an un existing file " + m.fullName)
	}
	m.fs.fileData[m.fullName] = m.data
	return nil
}

func (m *MemFile) Read(p []byte) (n int, err error) {
	n = copy(p, m.data[m.local:])
	m.local += n
	return n, nil
}

func (m *MemFile) ReadString(limit int) (string, error) {
	data := make([]byte, limit)
	n, err := m.Read(data)
	if err != nil {
		return "", err
	}
	return string(data[:n]), nil
}

func (m *MemFile) Write(p []byte) (n int, err error) {
	wLen := len(p)
	if m.local+wLen > len(m.data) {
		// 需要新开空间
		m.data = append(m.data[0:m.local], p...)
	} else {
		// copy 即可
		copy(m.data[m.local:], p)
	}
	m.local += len(p)
	return len(p), nil
}

func (m *MemFile) Offset(offset int64) error {
	if offset < 0 {
		return errors.New("File " + m.fullName + " cannot offset" + strconv.Itoa(int(offset)))
	}
	if int(offset) > len(m.data) {
		logger.Info("File original length", len(m.data), "seeks", offset)
		m.data = append(m.data, make([]byte, int(offset)-len(m.data))...)
		logger.Info("File new length", len(m.data))
	}
	m.local = int(offset)
	return nil
}

func (m *MemFileMeta) BaseName() string {
	return m.baseName
}

func (m *MemFileMeta) FullName() string {
	return m.fullName
}

func (m *MemFileMeta) Size() int64 {
	return m.size
}

func (m *MemFileMeta) IsDirectory() bool {
	return m.isDirectory
}

func (m *MemFileMeta) String() string {
	return fmt.Sprintf("%s[%s] %d bytes dir[%v]", m.baseName, m.fullName, m.size, m.isDirectory)
}
