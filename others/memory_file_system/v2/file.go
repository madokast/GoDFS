package v2

import (
	"errors"
	"fmt"
	"github.com/madokast/GoDFS/utils/logger"
	"strconv"
	"strings"
)

/**
FileIO 因为有一个读写位置的概念，所以都要加锁
*/

/*============================= FileIO ==========================================*/

func (m *MemFileDescription2) Read(p []byte) (n int, err error) {
	// 第一个锁，为了记录偏移
	m.fdLock.Lock()
	defer m.fdLock.Unlock()
	// 第二个锁，为了读取 openMode
	m.meta.openModeLock.RLock()
	defer m.meta.openModeLock.RUnlock()
	if !canReadMode(m.meta.openMode.Load()) {
		return 0, errors.New("Cannot read " + m.meta.fullName + " because not on read mode. Mode " + modeString(m.meta.openMode.Load()))
	}
	n = copy(p, m.meta.data[m.local:])
	m.local += int64(n)
	return n, nil
}

func (m *MemFileDescription2) Write(p []byte) (n int, err error) {
	m.fdLock.Lock()
	defer m.fdLock.Unlock()
	m.meta.openModeLock.RLock()
	defer m.meta.openModeLock.RUnlock()
	if m.meta.openMode.Load() != writeMode {
		return 0, errors.New("Cannot write " + m.meta.fullName + " because not on write mode")
	}
	n = copy(m.meta.data[m.local:], p)
	m.local += int64(n)
	return n, nil
}

func (m *MemFileDescription2) Offset(offset int64) error {
	m.fdLock.Lock()
	defer m.fdLock.Unlock()
	m.meta.openModeLock.RLock()
	defer m.meta.openModeLock.RUnlock()
	if m.meta.openMode.Load() == closeMode {
		return errors.New("Cannot offset " + m.meta.fullName + " because file is closed.")
	}
	if offset < 0 || offset >= m.meta.size {
		return errors.New(fmt.Sprintf("offset %d invaid for file %s size %d", offset, m.meta.fullName, m.meta.size))
	}
	m.local = offset
	return nil
}

// ReadString impl from Read
func (m *MemFileDescription2) ReadString(limit int) (string, error) {
	data := make([]byte, limit)
	n, err := m.Read(data)
	if err != nil {
		return "", err
	}
	return string(data[:n]), nil
}

func (m *MemFileDescription2) Close() error {
	m.fdLock.Lock()
	defer m.fdLock.Unlock()
	if m.closed {
		logger.Warn("Do not close file", m.meta.fullName, "more than once.")
		return nil
	}

	m.meta.openModeLock.Lock()
	defer m.meta.openModeLock.Unlock()

	curMode := m.meta.openMode.Load()
	if curMode == closeMode {
		panic("???")
	}
	newMode := closeMode
	if isReadMode(curMode) {
		readRef := curMode >> 16
		logger.Debug("Before close", m.meta.fullName, modeString(curMode))
		readRef--
		if readRef > 0 {
			newMode = (readRef << 16) | readMode
		}
	}

	if !m.meta.openMode.CompareAndSwap(curMode, newMode) {
		panic("???")
	}
	m.closed = true
	logger.Debug("After close", m.meta.fullName, modeString(newMode))
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

func (m *MemFileMeta2) String() string {
	sb := strings.Builder{}
	sb.WriteString(m.fullName)
	if m.isDirectory {
		sb.WriteString("[dir]")
		sb.WriteString(" child[" + strconv.Itoa(len(m.objects)) + "]")
	} else {
		sb.WriteString("[" + strconv.FormatInt(m.size, 10) + "B]")

		m.openModeLock.RLock()
		defer m.openModeLock.RUnlock()
		sb.WriteString(" " + modeString(m.openMode.Load()))
	}
	return sb.String()
}
