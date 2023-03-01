package jsonutils

import (
	"encoding/json"
	"io"
)

func Unmarshal(reader io.Reader, val interface{}) error {
	//data := make([]byte, 1024)
	//n, err2 := reader.Read(data)
	//logger.Info(n, err2, string(data[:n]))
	//os.Exit(1)

	decoder := json.NewDecoder(reader)
	err := decoder.Decode(val)
	return err
}

func Marshal(writer io.Writer, val interface{}) error {
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(val)
	return err
}

func String(val interface{}) string {
	marshal, err := json.Marshal(val)
	if err != nil {
		return err.Error()
	}
	return string(marshal)
}
