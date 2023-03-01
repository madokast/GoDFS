package example_node

import (
	"errors"
	"fmt"
	"github.com/madokast/GoDFS/internal/web"
	"github.com/madokast/GoDFS/utils/httputils"
	"github.com/madokast/GoDFS/utils/logger"
	"net/http"
)

type node struct {
	IP   string `json:"ip,omitempty"`
	Port uint16 `json:"port,omitempty"`
	Info string `json:"info,omitempty"`
}

func (n *node) serverGo() {
	http.HandleFunc("/info", func(w http.ResponseWriter, r *http.Request) {
		logger.Info("listen /info")
		httputils.Handle(w, r, &web.NullRequest{}, func(req *web.NullRequest) *node {
			ret := &node{IP: n.IP, Port: n.Port, Info: n.Info + " ^_^"}
			return ret
		})
	})

	go func() {
		err := http.ListenAndServe(n.addr(), nil)
		if err != nil {
			logger.Error(err)
		}
	}()
}

func (n *node) getInfo() (*node, error) {
	ret := &web.Response[*node]{}
	logger.Info("call /info")
	err := httputils.Post(n.IP, n.Port, "/info", &web.NullRequest{}, ret)
	if ret.Msg != web.SuccessMsg {
		return nil, errors.New(ret.Msg)
	}
	return ret.Data, err
}

func (n *node) addr() string {
	return fmt.Sprintf("%s:%d", n.IP, n.Port)
}
