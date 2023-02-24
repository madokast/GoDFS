package file_download

import (
	"bytes"
	"github.com/madokast/GoDFS/utils/logger"
	"io"
	"net/http"
	"strconv"
)

/*
测试文件下载请求
*/

func runServer(addr string) {
	// 假设存在文件 1.txt，内容为
	fileName := "1.txt"
	fileData := []byte("好好学习，天天向上。")

	http.HandleFunc("/1.txt", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Disposition", "attachment; filename="+fileName)
		w.Header().Set("Content-Type", http.DetectContentType(fileData)) // text/plain; charset=utf-8
		w.Header().Set("Content-Length", strconv.Itoa(len(fileData)))

		n, err := io.Copy(w, bytes.NewReader(fileData))
		logger.Info(n, err)
	})

	go func() {
		err := http.ListenAndServe(addr, nil)
		if err != nil {
			logger.Error(err)
		}
	}()
}
