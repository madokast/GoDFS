package web

import "net/http"

// HandleFunc http 请求处理函数
type HandleFunc = func(http.ResponseWriter, *http.Request)
