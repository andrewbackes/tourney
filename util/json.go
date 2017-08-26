package util

import (
	"encoding/json"
	"io"
)

func WriteJSON(obj interface{}, w io.Writer) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(obj)
}

func ReadJSON(reader io.Reader, dest interface{}) {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&dest)
	if err != nil {
		panic(err)
	}
}
