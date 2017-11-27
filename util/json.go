package util

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

func WriteJSON(obj interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
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

func PostJSON(url string, obj interface{}) {
	jsonValue, _ := json.Marshal(obj)
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		panic(err)
	}
	if resp.StatusCode > 399 {
		panic(resp.StatusCode)
	}
}
