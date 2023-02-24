package http_proxy

import (
	"github.com/madokast/GoDFS/utils/logger"
	"net/http"
	"strings"
	"testing"
	"time"
)

func Test_runServers(t *testing.T) {
	address := []string{
		"localhost:8081",
		"localhost:8082",
		"localhost:8083",
	}
	runServers(address)

	time.Sleep(50 * time.Millisecond)

	for i := 0; i < 10; i++ {
		addr := address[i%len(address)]
		url := "http" + "://" + addr + "/proxy"
		resp, err := http.Get(url)
		if err != nil && !strings.Contains(err.Error(), "EOF") {
			logger.Error(err)
		} else {
			bytes := make([]byte, 128, 128)
			n, _ := resp.Body.Read(bytes)
			logger.Info("resp", string(bytes[:n]))
		}
	}

}
