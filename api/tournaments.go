package api

import (
	"fmt"
	"github.com/andrewbackes/tourney/model"
	"github.com/andrewbackes/tourney/model/structures"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func (c *controller) getTournaments(w http.ResponseWriter, req *http.Request) {
	t := c.model.GetTournaments()
	writeJSON(t, w)
}

func (c *controller) getTournament(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := bson.ObjectIdHex(vars["tid"])
	t, err := c.model.GetTournament(id)
	if err == nil {
		writeJSON(t, w)
	} else if err == model.ErrorNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else {
		panic(err)
	}
}

func (c *controller) postTournament(w http.ResponseWriter, req *http.Request) {
	var t structures.Tournament
	readJSON(req.Body, &t)
	defer req.Body.Close()
	id, err := c.model.AddTournament(&t)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		panic(err)
	}
	w.Write([]byte(fmt.Sprintf("{\"id\":\"%s\"}", id.Hex())))
}

func (c *controller) getTournamentsNextGame(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	id := bson.ObjectIdHex(vars["tid"])
	t, _ := c.model.GetTournament(id)
	g := t.NextGame()
	if g != nil {
		writeJSON(g, w)
	} else {
		w.Write([]byte("{}"))
	}
}
