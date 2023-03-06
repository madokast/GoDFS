package fs

type Meta interface {
	FullName() string
	Exist() bool
	Size() int64
	IsDirectory() bool
	ModifyTime() int64
	String() string
}
