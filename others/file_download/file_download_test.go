package file_download

import (
	"github.com/madokast/GoDFS/utils/logger"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_runServer(t *testing.T) {
	runServer(":8080")

	time.Sleep(50 * time.Millisecond)

	resp, err := http.Get("http://localhost:8080/1.txt")
	if err != nil {
		logger.Error(err)
	}
	data := make([]byte, 1024)
	n, err := resp.Body.Read(data)
	if err != nil && !strings.Contains(err.Error(), "EOF") {
		logger.Error(err)
	} else {
		logger.Info("1.txt", string(data[:n]))
	}
}
