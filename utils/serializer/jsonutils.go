package serializer

import (
	"encoding/gob"
	"encoding/json"
	"io"
)

func JsonUnmarshal(reader io.Reader, val interface{}) error {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(val)
	return err
}

func JsonMarshal(writer io.Writer, val interface{}) error {
	encoder := json.NewEncoder(writer)
	err := encoder.Encode(val)
	return err
}

func GobUnmarshal(reader io.Reader, val interface{}) error {
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(val)
	return err
}

func GobMarshal(writer io.Writer, val interface{}) error {
	encoder := gob.NewEncoder(writer)
	err := encoder.Encode(val)
	return err
}

func JsonString(val interface{}) string {
	marshal, err := json.Marshal(val)
	if err != nil {
		return err.Error()
	}
	return string(marshal)
}
