package httputils

import (
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/logger"
	"github.com/madokast/GoDFS/utils/serializer"
	"net/http"
	"strconv"
	"testing"
	"time"
)

type person struct {
	Name string `json:"name,omitempty"`
}

type hello struct {
	Name string `json:"name,omitempty"`
	Time int64  `json:"time,omitempty"`
}

func TestPost(t *testing.T) {

	http.HandleFunc("/hello", func(writer http.ResponseWriter, request *http.Request) {
		p := person{}
		HandleJson(writer, request, &p, func(req *person) (*hello, error) {
			logger.Info(p)
			return &hello{Name: p.Name, Time: time.Now().UnixMilli()}, nil
		})
	})

	freePort := GetFreePort()

	go func() {
		err := http.ListenAndServe("127.0.0.1:"+strconv.Itoa(int(freePort)), nil)
		if err != nil {
			logger.Error(err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	p := person{Name: "xyn"}
	h := web.Response[*hello]{}
	err := PostJson("127.0.0.1", freePort, "/hello", &p, &h)
	if err != nil && err.Error() != "EOF" {
		logger.Error(err)
	}
	logger.Info(serializer.JsonString(h))
}

func TestGetFreePort(t *testing.T) {
	logger.Info(GetFreePort())
	logger.Info(GetFreePort())
	logger.Info(GetFreePort())
	logger.Info(GetFreePort())
}
