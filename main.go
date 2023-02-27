package main

import (
	"encoding/json"
	"github.com/madokast/GoDFS/utils/logger"
	"io"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/time", func(w http.ResponseWriter, r *http.Request) {
		// 允许跨域
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST,GET,PUT,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "content-type")

		if r.Method == http.MethodPost {
			data, err := io.ReadAll(r.Body)
			if err != nil {
				logger.Error(err)
			}
			logger.Info("data", string(data))
			var person struct {
				Name string `json:"name,omitempty"`
			}
			err = json.Unmarshal(data, &person)
			if err != nil {
				logger.Error(err)
			}
			var ret = struct {
				Name string `json:"name,omitempty"`
				Time int64  `json:"time,omitempty"`
			}{Name: person.Name, Time: time.Now().UnixMilli()}
			retBytes, err := json.Marshal(&ret)
			if err != nil {
				logger.Error(err)
			}
			w.Header().Set("content-type", "application/json")

			_, err = w.Write(retBytes)
			if err != nil {
				logger.Error(err)
			}
		}

	})

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.Error(err)
	}
}
