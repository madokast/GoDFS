package httputils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/logger"
	"github.com/madokast/GoDFS/utils/serializer"
	"net"
	"net/http"
)

// PostJson 发送 POST 请求，其中数据将序列化为 JSON
// req 是请求体，必须是指针类型
// resp 用返回结果，Response[*T] 类型，必须是指针类型，返回结果写入其中
func PostJson[Request, Response interface{}](ip string, port uint16, api string, req Request, resp Response) error {
	data, err := json.Marshal(req)
	if err != nil {
		return err
	}
	reader := bytes.NewReader(data)
	response, err := http.Post(fmt.Sprintf("http://%s:%d%s", ip, port, api), "application/json", reader)
	if err != nil {
		return err
	}

	err = serializer.JsonUnmarshal(response.Body, resp)
	return err
}

// PostGob 发送 POST 请求，其中数据将序列化为 Gob
func PostGob[Request, Response interface{}](ip string, port uint16, api string, req Request, resp Response) error {
	buf := bytes.Buffer{}
	err := serializer.GobMarshal(&buf, req)
	if err != nil {
		return err
	}

	response, err := http.Post(fmt.Sprintf("http://%s:%d%s", ip, port, api), "application/octet-stream", &buf)
	if err != nil {
		return err
	}

	err = serializer.GobUnmarshal(response.Body, resp)
	return err
}

// HandleJson 处理请求，其中数据以 JSON 发送收发
// writer 和 request 来自 http handler
// reqEmpty 是请求体，必须是指针，内容为空，用于存放请求信息，并传入处理逻辑 f 中处理
// f 是处理函数，传入请求体 Request，返回任意类型 T
// w 中写入包装后的 web.Response[*T] 类型，如果发生异常，返回 BaseResponse 类型
func HandleJson[Request, Response interface{}](writer http.ResponseWriter, request *http.Request, reqEmpty Request, f func(req Request) (Response, error)) {
	if request != interface{}(&web.NullRequest{}) {
		err := serializer.JsonUnmarshal(request.Body, reqEmpty)
		if err != nil {
			logger.Error(err)
			err = serializer.JsonMarshal(writer, web.Fail(404, "JsonUnmarshal request"))
			if err != nil {
				logger.Error(err)
			}
			return
		}
	}
	response, err := f(reqEmpty)
	if err != nil {
		err = serializer.JsonMarshal(writer, web.Fail(404, err.Error()))
		if err != nil {
			logger.Error(err)
		}
		return
	}

	err = serializer.JsonMarshal(writer, web.Success(response))
	if err != nil {
		logger.Error(err)
		err = serializer.JsonMarshal(writer, web.Fail(500, "JsonMarshal response"))
		if err != nil {
			logger.Error(err)
		}
	}
}

func HandleGob[Request, Response interface{}](writer http.ResponseWriter, request *http.Request, reqEmpty Request, f func(req Request) (Response, error)) {
	if request != interface{}(&web.NullRequest{}) {
		err := serializer.GobUnmarshal(request.Body, reqEmpty)
		if err != nil {
			logger.Error(err)
			err = serializer.GobMarshal(writer, web.Fail(404, "GobUnmarshal request"))
			if err != nil {
				logger.Error(err)
			}
			return
		}
	}
	response, err := f(reqEmpty)
	if err != nil {
		err = serializer.GobMarshal(writer, web.Fail(404, err.Error()))
		if err != nil {
			logger.Error(err)
		}
		return
	}

	err = serializer.GobMarshal(writer, web.Success(response))
	if err != nil {
		logger.Error(err)
		err = serializer.GobMarshal(writer, web.Fail(500, "GobMarshal response"))
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
