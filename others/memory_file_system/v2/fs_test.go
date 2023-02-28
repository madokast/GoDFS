package v2

import (
	"github.com/madokast/GoDFS/utils"
	"github.com/madokast/GoDFS/utils/logger"
	"testing"
)

func TestMemFileSystem2_ListAllPath(t *testing.T) {
	fs := NewMemFS()
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_Stat(t *testing.T) {
	fs := NewMemFS()
	logger.Info(fs.Stat("/"))
}

func TestMemFileSystem2_MakeDirectory(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectory("/a")
	if err != nil {
		panic(err)
	}
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_MakeDirectory2(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectory("/a")
	if err != nil {
		panic(err)
	}
	err = fs.MakeDirectory("/b")
	if err != nil {
		panic(err)
	}
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_MakeDirectory3(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectory("/a")
	if err != nil {
		panic(err)
	}
	err = fs.MakeDirectory("/a/b")
	if err != nil {
		panic(err)
	}
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func Test_fatherPath(t *testing.T) {
	logger.Info(fatherPath("/"))
	logger.Info(fatherPath("/a"))
	logger.Info(fatherPath("/a/b"))
	logger.Info(fatherPath("/a/b/1.txt"))
}

func Test_split(t *testing.T) {
	logger.Info(hierarchyFullPaths("/"))
	logger.Info(hierarchyFullPaths("/aa"))
	logger.Info(hierarchyFullPaths("/aa/bb"))
	logger.Info(hierarchyFullPaths("/aa/bb/2.txt"))
}

func TestMemFileSystem2_CreateFile(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/1.txt", 4)
	utils.PanicIfErr(err)
	fileIO, err := fs.OpenFile("/1.txt", true)
	utils.PanicIfErr(err)
	n, err := fileIO.Write([]byte("abcd"))
	utils.PanicIfErr(err)
	utils.PanicIf(n != 4, n)
	err = fileIO.Offset(0)
	utils.PanicIfErr(err)
	readString, err := fileIO.ReadString(4)
	utils.PanicIfErr(err)
	utils.PanicIf(readString != "abcd", readString)
	err = fileIO.Close()
	utils.PanicIfErr(err)
	exist := fs.Exist("/1.txt")
	utils.PanicIf(!exist)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta)
	}
}

func TestMemFileSystem2_OpenFile(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/1.txt", 4)
	utils.PanicIfErr(err)
	fileIO, err := fs.OpenFile("/1.txt", true)
	utils.PanicIfErr(err)
	err = fs.DeleteFile("/1.txt")
	logger.Info(err)
	utils.PanicIf(err == nil, err)

	err = fileIO.Close()
	utils.PanicIfErr(err)
	exist := fs.Exist("/1.txt")
	utils.PanicIf(!exist)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta)
	}

	err = fs.DeleteFile("/1.txt")
	utils.PanicIfErr(err)
}

func TestMemFileSystem2_OpenFile2(t *testing.T) {
	fs := NewMemFS()
	_, err := fs.OpenFile("/1.txt", true)
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_OpenFile3(t *testing.T) {
	fs := NewMemFS()
	_, err := fs.OpenFile("1.txt", true)
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_CreateFile1(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("//1.txt", 4)
	logger.Info(err)
	utils.PanicIf(err == nil, nil)
}

func TestMemFileSystem2_CreateFile2(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/", 4)
	logger.Info(err)
	utils.PanicIf(err == nil, nil)
}

func TestMemFileSystem2_CreateFile3(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/1.txt", 4)
	utils.PanicIfErr(err)
	err = fs.CreateFile("/1.txt/2.txt", 4)
	logger.Info(err)
	utils.PanicIf(err == nil, nil)
}

func TestMemFileSystem2_CreateFile4(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/1.txt", 4)
	utils.PanicIfErr(err)
	err = fs.CreateFile("/1.txt", 5)
	logger.Info(err)
	utils.PanicIf(err == nil, nil)
}

func TestMemFileSystem2_OpenFile4(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/a.txt", 4)
	utils.PanicIfErr(err)
	_, err = fs.OpenFile("/a.txt", true)
	utils.PanicIfErr(err)
	_, err = fs.OpenFile("/a.txt", true)
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_OpenFile5(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/a.txt", 4)
	utils.PanicIfErr(err)
	_, err = fs.OpenFile("/a.txt", true)
	utils.PanicIfErr(err)
	_, err = fs.OpenFile("/a.txt", false)
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_OpenFile6(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/a.txt", 4)
	utils.PanicIfErr(err)
	_, err = fs.OpenFile("/a.txt", false)
	utils.PanicIfErr(err)
	_, err = fs.OpenFile("/a.txt", true)
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_OpenFile7(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/a.txt", 4)
	utils.PanicIfErr(err)
	stat, err := fs.Stat("/a.txt")
	utils.PanicIfErr(err)
	io1, err := fs.OpenFile("/a.txt", false)
	utils.PanicIfErr(err)
	logger.Info(stat)
	io2, err := fs.OpenFile("/a.txt", false)
	utils.PanicIfErr(err)
	logger.Info(stat)
	io3, err := fs.OpenFile("/a.txt", false)
	utils.PanicIfErr(err)
	logger.Info(stat)

	err = io2.Close()
	utils.PanicIfErr(err)
	logger.Info(stat)

	logger.Info("=== warn ok ===")
	err = io2.Close()
	utils.PanicIfErr(err)
	logger.Info(stat)

	logger.Info("=== warn ok ===")
	err = io2.Close()
	utils.PanicIfErr(err)
	logger.Info(stat)

	err = io3.Close()
	utils.PanicIfErr(err)
	logger.Info(stat)

	err = io1.Close()
	utils.PanicIfErr(err)
	logger.Info(stat)
}

func TestMemFileSystem2_DeleteFile(t *testing.T) {
	fs := NewMemFS()
	err := fs.DeleteFile("/")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_DeleteFile2(t *testing.T) {
	fs := NewMemFS()
	err := fs.DeleteFile("1.txt")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_DeleteFile3(t *testing.T) {
	fs := NewMemFS()
	logger.Info("==== warn ok ====")
	err := fs.DeleteFile("/1.txt")
	utils.PanicIfErr(err)
}

func TestMemFileSystem2_DeleteFile4(t *testing.T) {
	fs := NewMemFS()
	logger.Info("==== warn ok ====")
	err := fs.DeleteFile("/a/b/c/1.txt")
	utils.PanicIfErr(err)
}

func TestMemFileSystem2_DeleteFile5(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a")
	utils.PanicIfErr(err)
	err = fs.DeleteFile("/a")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_DeleteFile6(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/a", 10)
	utils.PanicIfErr(err)
	err = fs.DeleteFile("/a")
	//logger.Info(err)
	utils.PanicIfErr(err)
}

func TestMemFileSystem2_DeleteFile7(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/a", 10)
	utils.PanicIfErr(err)
	_, err = fs.OpenFile("/a", false)
	utils.PanicIfErr(err)
	err = fs.DeleteFile("/a")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_DeleteFile8(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/a", 10)
	utils.PanicIfErr(err)
	_, err = fs.OpenFile("/a", true)
	utils.PanicIfErr(err)
	err = fs.DeleteFile("/a")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_DeleteFile9(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/a", 10)
	utils.PanicIfErr(err)
	_, err = fs.OpenFile("/a", false)
	utils.PanicIfErr(err)
	_, err = fs.OpenFile("/a", false)
	utils.PanicIfErr(err)
	err = fs.DeleteFile("/a")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_MakeDirectory1(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectory("/a/b")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_MakeDirectory4(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectory("adfaf")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_MakeDirectory5(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectory("/")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_MakeDirectory6(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectory("/a")
	utils.PanicIfErr(err)
	err = fs.MakeDirectory("/a")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_MakeDirectory7(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/1.txt", 10)
	utils.PanicIfErr(err)
	err = fs.MakeDirectory("/1.txt")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
}

func TestMemFileSystem2_MakeDirectories(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a/b")
	utils.PanicIfErr(err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_DeleteDirectory(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a/b")
	utils.PanicIfErr(err)
	err = fs.DeleteDirectory("/a/b")
	utils.PanicIfErr(err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_DeleteDirectory1(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a/b")
	utils.PanicIfErr(err)
	err = fs.DeleteDirectory("/a")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_DeleteDirectory2(t *testing.T) {
	fs := NewMemFS()
	err := fs.CreateFile("/a", 1)
	utils.PanicIfErr(err)
	err = fs.DeleteDirectory("/a")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_DeleteDirectory3(t *testing.T) {
	fs := NewMemFS()
	logger.Info("=== warn ok ===")
	err := fs.DeleteDirectory("/a")
	utils.PanicIfErr(err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_DeleteDirectoryAll(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a/b")
	utils.PanicIfErr(err)
	err = fs.DeleteDirectoryAll("/a")
	utils.PanicIfErr(err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_DeleteDirectoryAll2(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a/b/c")
	utils.PanicIfErr(err)

	err = fs.MakeDirectories("/a/f/g")
	utils.PanicIfErr(err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}

	err = fs.DeleteDirectoryAll("/a")
	utils.PanicIfErr(err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_DeleteDirectoryAll3(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a/b")
	utils.PanicIfErr(err)

	utils.PanicIfErr(fs.CreateFile("/a/1.txt", 1))
	utils.PanicIfErr(fs.CreateFile("/a/b/1.txt", 2))

	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}

	err = fs.DeleteDirectoryAll("/a")
	utils.PanicIfErr(err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_DeleteDirectoryAll4(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a/b")
	utils.PanicIfErr(err)

	utils.PanicIfErr(fs.CreateFile("/a/1.txt", 1))
	utils.PanicIfErr(fs.CreateFile("/a/b/2.txt", 2))

	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}

	_, err = fs.OpenFile("/a/b/2.txt", false)
	utils.PanicIfErr(err)

	err = fs.DeleteDirectoryAll("/a")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_DeleteDirectoryAll5(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a/b")
	utils.PanicIfErr(err)

	utils.PanicIfErr(fs.CreateFile("/a/1.txt", 1))
	utils.PanicIfErr(fs.CreateFile("/a/b/2.txt", 2))

	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}

	io, err := fs.OpenFile("/a/b/2.txt", false)
	utils.PanicIfErr(err)

	err = fs.DeleteDirectoryAll("/a")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}

	utils.PanicIfErr(io.Close())

	err = fs.DeleteDirectoryAll("/a")
	utils.PanicIfErr(err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_DeleteDirectoryAll6(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a/b")
	utils.PanicIfErr(err)

	utils.PanicIfErr(fs.CreateFile("/a/1.txt", 1))
	utils.PanicIfErr(fs.CreateFile("/a/b/2.txt", 2))

	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}

	io, err := fs.OpenFile("/a/b/2.txt", true)
	utils.PanicIfErr(err)

	err = fs.DeleteDirectoryAll("/a")
	logger.Info(err)
	utils.PanicIf(err == nil, err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}

	utils.PanicIfErr(io.Close())

	err = fs.DeleteDirectoryAll("/a")
	utils.PanicIfErr(err)
	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}
}

func TestMemFileSystem2_ReadDirectory(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a/b")
	utils.PanicIfErr(err)

	utils.PanicIfErr(fs.CreateFile("/a/1.txt", 1))
	utils.PanicIfErr(fs.CreateFile("/a/b/2.txt", 2))

	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}

	files, dirs, err := fs.ReadDirectory("/a")
	utils.PanicIfErr(err)
	logger.Info("/a", "files", files)
	logger.Info("/a", "dirs", dirs)
}

func TestMemFileSystem2_ReadDirectory2(t *testing.T) {
	fs := NewMemFS()
	err := fs.MakeDirectories("/a/b")
	utils.PanicIfErr(err)

	utils.PanicIfErr(fs.CreateFile("/a/1.txt", 1))
	utils.PanicIfErr(fs.CreateFile("/a/b/2.txt", 2))

	for _, meta := range fs.ListAllPath() {
		logger.Info(meta.FullName())
	}

	files, dirs, err := fs.ReadDirectory("/a/b")
	utils.PanicIfErr(err)
	logger.Info("/a/b", "files", files)
	logger.Info("/a/b", "dirs", dirs)
}
