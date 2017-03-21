package api

import (
	"encoding/json"
	"fmt"
	"github.com/andrewbackes/tourney/model"
	"github.com/andrewbackes/tourney/model/structures"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io"
	"net/http"
)

type controller struct {
	model *model.Model
}

func writeJSON(obj interface{}, w io.Writer) {
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(obj)
}

func (c *controller) getTournaments(w http.ResponseWriter, req *http.Request) {
	t := c.model.GetTournaments()
	writeJSON(t, w)
}

func (c *controller) getTournament(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := bson.ObjectIdHex(vars["id"])
	t, err := c.model.GetTournament(id)
	if err == nil {
		writeJSON(t, w)
	} else if err.Error() == "Tournament not found" {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		panic(err)
	}
}

func (c *controller) postTournaments(w http.ResponseWriter, req *http.Request) {
	decoder := json.NewDecoder(req.Body)
	var t structures.Tournament
	err := decoder.Decode(&t)
	if err != nil {
		panic(err)
	}
	defer req.Body.Close()
	id := c.model.CreateTournament(t)
	w.Write([]byte(fmt.Sprintf("{\"id\":\"%s\"}", id.Hex())))
}
