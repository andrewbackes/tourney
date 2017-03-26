package api

import (
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func (c *controller) getGames(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tid := bson.ObjectIdHex(vars["tid"])
	t, _ := c.model.GetTournament(tid)
	g := t.GetGames()
	writeJSON(g, w)
}

func (c *controller) getGame(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tid := bson.ObjectIdHex(vars["tid"])
	gid := bson.ObjectIdHex(vars["gid"])
	t, _ := c.model.GetTournament(tid)
	g := t.GetGame(gid)
	writeJSON(g, w)
}
