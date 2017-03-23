package api

import (
	"encoding/json"
	"github.com/andrewbackes/tourney/model"
	"io"
)

type controller struct {
	model *model.Model
}

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
