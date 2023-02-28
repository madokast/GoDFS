package v2

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/madokast/GoDFS/internal/ifile"
	"github.com/madokast/GoDFS/utils/logger"
	"sort"
	"strings"
	"time"
)

/**
基于内存的文件系统
所有路径名必须是全称 /a/b/c 形式
*/

// CreateFile 创建文件，文件大小必须提前指定，不能改变。返回打开的文件 IO
// 例如创建文件 /a/b/1.txt
// 首先读取路径 /a/b 确认存在，并确认其是目录
// 分配 size 空间，创建文件元信息
// 在目录 /a/b 中写入文件全称 /a/b/1.txt
// 将文件元信息写入大 map 中
// 最后调用 OpenFile 方法
func (fs *MemFileSystem2) CreateFile(name string, size int64) error {
	err := checkPath(name)
	if err != nil {
		return errors.New("Cannot create file because " + err.Error())
	}

	father, exist := fatherPath(name)
	if !exist {
		return errors.New("Cannot create file " + name)
	}

	/*-------------------------------------- FS 写临界区 -----------------------------------------*/

	fs.fsLock.Lock()
	defer fs.fsLock.Unlock()
	ex := fs.existUnlock(name)
	if ex {
		return errors.New("Cannot create file " + name + " because it exists.")
	}

	fatherMeta, err := fs.statUnlock(father)
	if err != nil {
		return err
	}

	if !fatherMeta.IsDirectory() {
		return errors.New("Cannot create " + name + " because path " + father + " is not a dir.")
	}

	fatherMeta2 := fatherMeta.(*MemFileMeta2)
	_, ex = fatherMeta2.objects[name]
	if ex {
		logger.Error(name, fatherMeta2.objects)
		panic("???")
	}
	fatherMeta2.objects[name] = struct{}{}
	fs.fileMap[name] = &MemFileMeta2{
		fullName:    name,
		size:        size,
		isDirectory: false,
		data:        make([]byte, size),
		objects:     nil,
	}
	return nil
}

// OpenFile 打开文件
// 直接去大 map 中找到文件元信息
// 确定文件存在，是文件不是目录
// 从元信息中拿到底层数据块，包装为 IO
func (fs *MemFileSystem2) OpenFile(name string, write bool) (ifile.FileIO, error) {
	err := checkPath(name)
	if err != nil {
		return nil, errors.New("Cannot open file because " + err.Error())
	}

	/*-------------------------------------- FS 读临界区 -----------------------------------------*/

	fs.fsLock.Lock()
	defer fs.fsLock.Unlock()

	meta, err := fs.statUnlock(name)
	if err != nil {
		return nil, err
	}

	if meta.IsDirectory() {
		return nil, errors.New(name + " is a dir.")
	}

	meta2 := meta.(*MemFileMeta2)

	meta2.openModeLock.Lock()
	defer meta2.openModeLock.Unlock()

	var swapped bool
	if write {
		// 只允许从 close 到写模式
		swapped = meta2.openMode.CompareAndSwap(closeMode, writeMode)
	} else {
		// 尝试从 close 到读模式。1 记录打开次数 1
		swapped = meta2.openMode.CompareAndSwap(closeMode, (1<<16)|readMode)
		if !swapped {
			// 失败了则增加读引用数目。不需要自旋，因为所有 mode 修改都在 fs 全局锁下
			curMode := meta2.openMode.Load()
			if isReadMode(curMode) {
				readRef := curMode >> 16
				logger.Debug("Before read open", name, modeString(curMode))
				readRef++
				newMode := (readRef << 16) | readMode
				swapped = meta2.openMode.CompareAndSwap(curMode, newMode)
				if !swapped {
					panic("???")
				}
				logger.Debug("After read open", name, modeString(newMode))
			}
		}
	}
	if !swapped {
		curMode := meta2.openMode.Load()
		return nil, errors.New("File " + name + " has been opened. Mode " + modeString(curMode))
	}

	return &MemFileDescription2{
		local: 0,
		meta:  meta2,
	}, nil
}

// DeleteFile 删除文件，文件不存在只警告
func (fs *MemFileSystem2) DeleteFile(name string) error {
	return fs.deleteFile0(name, true)
}

func (fs *MemFileSystem2) deleteFile0(name string, lock bool) error {
	err := checkPath(name)
	if err != nil {
		return errors.New("Cannot delete file " + name + " because " + err.Error())
	}

	father, ok := fatherPath(name)
	if !ok {
		return errors.New("Cannot delete file " + name + " which is not a file.")
	}

	/*-------------------------------------- FS 写临界区 -----------------------------------------*/
	if lock {
		fs.fsLock.Lock()
		defer fs.fsLock.Unlock()
	}
	exist := fs.existUnlock(name)
	if !exist {
		logger.Warn("Cannot delete file " + name + " which does not exist.")
		return nil
	}
	fileMeta, err := fs.statUnlock(name)
	if err != nil {
		return err
	}
	if fileMeta.IsDirectory() {
		return errors.New("Cannot delete file " + name + " which is a dir.")
	}
	fileMeta2 := fileMeta.(*MemFileMeta2)

	fileMeta2.openModeLock.RLock()
	defer fileMeta2.openModeLock.RUnlock()
	if fileMeta2.openMode.Load() != closeMode {
		return errors.New(name + " is open. Mode " + modeString(fileMeta2.openMode.Load()))
	}
	fatherMeta, err := fs.statUnlock(father)
	if err != nil {
		return err
	}
	fatherMeta2 := fatherMeta.(*MemFileMeta2)
	if !fatherMeta2.isDirectory {
		logger.Error(name, fatherMeta2)
		panic("???")
	}
	_, ok = fatherMeta2.objects[name]
	if !ok {
		logger.Error(name, fatherMeta2.objects)
		panic("???")
	}
	delete(fatherMeta2.objects, name)
	// ... 假装释放底层资源
	time.Sleep(10 * time.Millisecond)
	delete(fs.fileMap, name)
	return nil
}

func (fs *MemFileSystem2) MakeDirectory(name string) error {
	return fs.makeDirectory0(name, true)
}

func (fs *MemFileSystem2) makeDirectory0(name string, lock bool) error {
	err := checkPath(name)
	if err != nil {
		return err
	}

	father, exist := fatherPath(name)
	if !exist {
		return errors.New("Cannot mkdir root path " + name + ".")
	}

	/*-------------------------------------- FS 写临界区 -----------------------------------------*/
	if lock {
		fs.fsLock.Lock()
		defer fs.fsLock.Unlock()
	}

	ex := fs.existUnlock(name)
	if ex {
		return errors.New("Cannot mkdir " + name + " which exists.")
	}

	fatherMeta, err := fs.statUnlock(father)
	if err != nil {
		return errors.New("Cannot mkdir " + name + " because father path " + err.Error())
	}
	if !fatherMeta.IsDirectory() {
		return errors.New("Cannot mkdir " + name + " because " + father + " is not a dir.")
	}

	fatherMeta2 := fatherMeta.(*MemFileMeta2)

	_, ex = fatherMeta2.objects[name]
	if ex {
		logger.Error(name, fatherMeta2.objects)
		panic("???")
	}
	fatherMeta2.objects[name] = struct{}{}
	fs.fileMap[name] = &MemFileMeta2{
		fullName:    name,
		size:        0,
		isDirectory: true,
		data:        nil,
		objects:     map[string]struct{}{},
	}
	return nil
}

func (fs *MemFileSystem2) MakeDirectories(name string) error {
	return fs.makeDirectories0(name, true)
}

func (fs *MemFileSystem2) makeDirectories0(name string, lock bool) error {
	err := checkPath(name)
	if err != nil {
		return errors.New("Cannot mkdir " + name + " because " + err.Error())
	}

	father, exist := fatherPath(name)
	if !exist {
		return errors.New("Cannot mkdir " + name + ".")
	}

	/*-------------------------------------- FS 写临界区 -----------------------------------------*/
	if lock {
		fs.fsLock.Lock()
		defer fs.fsLock.Unlock()
	}

	if fs.existUnlock(father) {
		// 父亲存在，则直接调用 makeDirectory0，已经在临界区不上锁。哎，Go 锁不可重入
		return fs.makeDirectory0(name, false)
	} else {
		// 父亲不存在，递归，已经在临界区不上锁
		err = fs.makeDirectories0(father, false)
		if err != nil {
			return err
		}
		return fs.makeDirectory0(name, false)
	}
}

func (fs *MemFileSystem2) DeleteDirectory(name string) error {
	return fs.deleteDirectory0(name, true)
}

func (fs *MemFileSystem2) deleteDirectory0(name string, lock bool) error {
	err := checkPath(name)
	if err != nil {
		return errors.New("Cannot delete dir " + name + " because " + err.Error())
	}
	father, ok := fatherPath(name)
	if !ok {
		return errors.New("cannot delete dir" + name)
	}
	/*-------------------------------------- FS 写临界区 -----------------------------------------*/
	if lock {
		fs.fsLock.Lock()
		defer fs.fsLock.Unlock()
	}
	exist := fs.existUnlock(name)
	if !exist {
		logger.Warn("cannot delete dir" + name + " which does not exist")
		return nil
	}
	fileMeta, err := fs.statUnlock(name)
	if err != nil {
		return err
	}
	if !fileMeta.IsDirectory() {
		return errors.New("cannot delete dir " + name + " which is a file")
	}
	fileMeta2 := fileMeta.(*MemFileMeta2)
	if len(fileMeta2.objects) != 0 {
		return errors.New("cannot delete dir" + name + " which is not empty")
	}
	fatherMeta, err := fs.statUnlock(father)
	if err != nil {
		return err
	}
	fatherFileMeta2 := fatherMeta.(*MemFileMeta2)
	if !fatherFileMeta2.isDirectory {
		logger.Error(name, fatherFileMeta2)
		panic("???")
	}
	_, ok = fatherFileMeta2.objects[name]
	if !ok {
		logger.Error(name, fatherFileMeta2)
		panic("???")
	}
	delete(fatherFileMeta2.objects, name)
	delete(fs.fileMap, name)
	return nil
}

func (fs *MemFileSystem2) DeleteDirectoryAll(name string) error {
	return fs.deleteDirectoryAll0(name, true)
}

func (fs *MemFileSystem2) deleteDirectoryAll0(name string, lock bool) error {
	err := checkPath(name)
	if err != nil {
		return errors.New("Cannot delete dir " + name + " because " + err.Error())
	}

	/*-------------------------------------- FS 写临界区 -----------------------------------------*/
	if lock {
		fs.fsLock.Lock()
		defer fs.fsLock.Unlock()
	}
	// 删除文件夹内内容
	exist := fs.existUnlock(name)
	if !exist {
		logger.Warn(name + " does not exist")
		return nil
	}
	fileMeta, err := fs.statUnlock(name)
	if err != nil {
		return err
	}
	if !fileMeta.IsDirectory() {
		return errors.New("Cannot delete dir " + name + " which is not a dir")
	}
	files, dirs, err := fs.readDirectory0(name, false)
	if err != nil {
		return errors.New("Cannot delete dir " + name + " because " + err.Error())
	}
	for _, file := range files {
		err = fs.deleteFile0(file, false)
		if err != nil {
			return errors.New("Cannot delete dir " + name + " because " + err.Error())
		}
	}
	for _, dir := range dirs {
		// 递归
		err = fs.deleteDirectoryAll0(dir, false)
		if err != nil {
			return errors.New("Cannot delete dir " + name + " because " + err.Error())
		}
	}
	// 最后删除自己
	return fs.deleteDirectory0(name, false)
}

func (fs *MemFileSystem2) ReadDirectory(dir string) (files []string, dirs []string, err error) {
	return fs.readDirectory0(dir, true)
}

func (fs *MemFileSystem2) readDirectory0(dir string, lock bool) (files []string, dirs []string, err error) {
	err = checkPath(dir)
	if err != nil {
		return nil, nil, err
	}
	/*-------------------------------------- FS 读临界区 -----------------------------------------*/
	if lock {
		fs.fsLock.Lock()
		defer fs.fsLock.Unlock()
	}
	fileMeta, err := fs.statUnlock(dir)
	if err != nil {
		return nil, nil, err
	}
	if !fileMeta.IsDirectory() {
		return nil, nil, errors.New(dir + " is not a dir")
	}
	fileMeta2 := fileMeta.(*MemFileMeta2)
	for k := range fileMeta2.objects {
		child, err := fs.statUnlock(k)
		if err != nil {
			panic(err)
		}
		if child.IsDirectory() {
			dirs = append(dirs, child.FullName())
		} else {
			files = append(files, child.FullName())
		}
	}
	return files, dirs, nil
}

func (fs *MemFileSystem2) Exist(name string) bool {
	fs.fsLock.Lock()
	defer fs.fsLock.Unlock()
	return fs.existUnlock(name)
}

func (fs *MemFileSystem2) existUnlock(name string) bool {
	_, ok := fs.fileMap[name]
	return ok
}

func (fs *MemFileSystem2) Stat(name string) (ifile.FileMeta, error) {
	err := checkPath(name)
	if err != nil {
		return nil, err
	}

	fs.fsLock.Lock()
	defer fs.fsLock.Unlock()
	return fs.statUnlock(name)
}

// statUnlock 不加锁的 Stat 供内部使用
func (fs *MemFileSystem2) statUnlock(name string) (ifile.FileMeta, error) {
	meta2, ok := fs.fileMap[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("%s does not exist.", name))
	}
	return meta2, nil
}

func (fs *MemFileSystem2) ListAllPath() []ifile.FileMeta {
	fs.fsLock.Lock()
	defer fs.fsLock.Unlock()

	var ret []ifile.FileMeta
	stack := list.New()
	stack.PushFront("/")
	for stack.Len() > 0 {
		pop := stack.Front()
		stack.Remove(pop)
		popMeta := fs.fileMap[pop.Value.(string)]
		ret = append(ret, popMeta)
		if popMeta.IsDirectory() {
			sortObjs := toSortedSlice(popMeta.objects)
			for i := len(sortObjs) - 1; i >= 0; i-- {
				stack.PushFront(sortObjs[i])
			}
		}
	}
	return ret
}

// hierarchyFullPaths 将合法的路径输出每个层级
// "/"            -> [/]
// "/aa"          -> [/ /aa]
// "/aa/bb"       -> [/ /aa /aa/bb]
// "/aa/bb/2.txt" -> [/ /aa /aa/bb /aa/bb/2.txt]
func hierarchyFullPaths(name string) []string {
	if name == "/" {
		return []string{"/"}
	}

	bases := strings.Split(name, "/")
	ret := []string{"/"}
	fullPath := ""
	for i := 1; i < len(bases); i++ {
		fullPath += "/" + bases[i]
		ret = append(ret, fullPath)
	}
	return ret
}

// checkPath 确认路径正确性
func checkPath(name string) error {
	// 必须 / 开头
	if !strings.HasPrefix(name, "/") {
		return errors.New("invalid path " + name + ". Path should start with /")
	}
	if name == "/" {
		return nil
	}
	// 非根目录，不能以 / 结尾
	if strings.HasSuffix(name, "/") {
		return errors.New("invalid path " + name + ". A non-root path should not end with /")
	}
	// 不能有连续的 //
	if strings.Contains(name, "//") {
		return errors.New("invalid path " + name + ". Path should not contain //")
	}
	return nil
}

// fatherPath 获取父路径。根目录没有父路径
func fatherPath(name string) (father string, exist bool) {
	// 定义根目录的父路径为 ""
	if name == "/" {
		return "", false
	}
	lastSlash := strings.LastIndex(name, "/")
	// 类似 /a 的一层路径，父路径为 /
	if lastSlash == 0 {
		return "/", true
	}
	// 类似 /a/b 的多层路径，父路径为 /a
	father = name[0:lastSlash]
	return father, true
}

func toSortedSlice(set map[string]struct{}) []string {
	r := make([]string, 0, len(set))
	for k := range set {
		r = append(r, k)
	}
	sort.Strings(r)
	return r
}
