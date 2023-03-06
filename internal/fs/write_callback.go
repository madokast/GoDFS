package fs

import "net/http"

/**
写通知回调
使用 RegisterWriteCallback 注册一个文件片段和回调函数
如果一个写操作涉及了该文件，并且写入范围和注册的范围有重叠，那么就会调用回调
回调是一次性的，调用后自动移除
可以手动删除注册的回调信息 RemoveWriteCallback
*/

type WriteCallBackObj struct {
	FileName string
	Offset   int64
	Length   int64
	Callback func()
}

type WriteCallBack interface {
	RegisterWriteCallback(*WriteCallBackObj)                   // 注册文件修改通知回调。缓存层需要用到，用来失效一些资源
	RemoveWriteCallback(*WriteCallBackObj)                     // 取消注册
	WriteCallback(fileName string, offset, length int64) error // 向 this 节点发送回调 WriteCallBack 请求
	DoWriteCallback(w http.ResponseWriter, r *http.Request)    // 处理 this 节点的回调 WriteCallBack 请求
}

// Intersect 范围是否相交
// 原理：发生相交，则覆盖范围肯定小于 length + obj.Length
func (obj *WriteCallBackObj) Intersect(offset, length int64) bool {
	end := offset + length - 1
	objEnd := obj.Offset + obj.Length - 1

	minStart := min(offset, obj.Offset)
	maxEnd := max(end, objEnd)

	return maxEnd-minStart+1 < (length + obj.Length)
}

func min(a, b int64) int64 {
	if a < b {
		return a
	} else {
		return b
	}
}

func max(a, b int64) int64 {
	if a > b {
		return a
	} else {
		return b
	}
}
