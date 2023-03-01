package http_proxy

import (
	"github.com/madokast/GoDFS/utils/logger"
	"math/rand"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type handle struct {
	selfAddress string
	addresses   []string
}

// 实现 http.Handler
func (h *handle) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	//logger.Info(path)
	if path == "/hello/proxy" {
		helloHandle(h)(w, r)
	} else if path == "/proxy" {
		helloProxy(h)(w, r)
	}
}

func helloHandle(h *handle) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("Do hello in " + h.selfAddress))
	}
}

// helloProxy 处理 /proxy 请求，随机转发到其他节点，url 变成 /hello/proxy
func helloProxy(h *handle) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		addr := h.addresses[rand.Intn(len(h.addresses))]
		parse, err := url.Parse("httputils" + "://" + addr + "/hello")
		if err != nil {
			logger.Error(err)
			_, _ = w.Write([]byte(err.Error()))
		} else {
			logger.Info(h.selfAddress, "proxy to", parse)
			proxy := httputil.NewSingleHostReverseProxy(parse)
			proxy.ServeHTTP(w, r)
		}
	}
}

func runServers(addresses []string) {
	for _, address := range addresses {
		go func(address string) {
			logger.Info("run", address)
			err := http.ListenAndServe(address, &handle{addresses: addresses, selfAddress: address})
			if err != nil {
				logger.Error(err)
			}
		}(address)
	}
}
