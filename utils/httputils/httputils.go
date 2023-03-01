package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/jsonutils"
	"github.com/madokast/GoDFS/utils/logger"
	"net"
	"net/http"
)

func Post[Request, Response interface{}](ip string, port uint16, api string, req Request, resp Response) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(data)
	response, err := http.Post(fmt.Sprintf("http://%s:%d%s", ip, port, api), "application/json", reader)
	if err != nil {
		return err
	}

	err = jsonutils.Unmarshal(response.Body, resp)
	return err
}

func Handle[Request, Response interface{}](writer http.ResponseWriter, request *http.Request, reqEmpty Request, f func(req Request) Response) {
	if request != interface{}(&web.NullRequest{}) {
		err := jsonutils.Unmarshal(request.Body, reqEmpty)
		if err != nil {
			logger.Error(err)
			err = jsonutils.Marshal(writer, web.Fail(404, "Unmarshal request"))
			if err != nil {
				logger.Error(err)
			}
			return
		}
	}
	response := f(reqEmpty)
	err := jsonutils.Marshal(writer, web.Success(response))
	if err != nil {
		logger.Error(err)
		err = jsonutils.Marshal(writer, web.Fail(500, "Marshal response"))
		if err != nil {
			logger.Error(err)
		}
	}
}

func GetFreePort() uint16 {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		logger.Error(err)
		return GetFreePort()
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		logger.Error(err)
		return GetFreePort()
	}
	defer func(listener *net.TCPListener) {
		err := listener.Close()
		if err != nil {
			logger.Error(err)
		}
	}(listener)
	return uint16(listener.Addr().(*net.TCPAddr).Port)
}
