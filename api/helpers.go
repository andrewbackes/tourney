package api

import (
	"encoding/json"
	"io"
)

func writeJSON(obj interface{}, w io.Writer) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(obj)
}

func readJSON(reader io.Reader, dest interface{}) {
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&dest)
	if err != nil {
		panic(err)
	}
}
