package web

import (
	"github.com/madokast/GoDFS/utils/serializer"
	"time"
)

const SuccessMsg = "Successful"

type NullRequest struct {
}

type NullResponse struct {
}

type BaseResponse struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	ServerTime int64  `json:"serverTime"`
}

type Response[T interface{}] struct {
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	ServerTime int64  `json:"serverTime"`
	Data       T      `json:"data"`
}

func Success[T interface{}](data T) *Response[T] {
	return &Response[T]{
		Code:       200,
		Msg:        SuccessMsg,
		ServerTime: time.Now().UnixMilli(),
		Data:       data,
	}
}

func Fail(code int, msg string) *BaseResponse {
	return &BaseResponse{
		Code:       code,
		Msg:        msg,
		ServerTime: time.Now().UnixMilli(),
	}
}

func (br *BaseResponse) String() string {
	return serializer.JsonString(br)
}
