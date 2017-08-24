package api

import (
	"fmt"
	"github.com/andrewbackes/tourney/helpers"
	"github.com/andrewbackes/tourney/model/structures"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func (c *controller) getGames(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tid := bson.ObjectIdHex(vars["tid"])
	t, _ := c.model.GetTournament(tid)
	g := t.GetGames()
	helpers.WriteJSON(g, w)
}

func (c *controller) getGame(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tid := bson.ObjectIdHex(vars["tid"])
	gid := bson.ObjectIdHex(vars["gid"])
	t, _ := c.model.GetTournament(tid)
	g := t.GetGame(gid)
	helpers.WriteJSON(g, w)
}

func (c *controller) patchGame(w http.ResponseWriter, req *http.Request) {
	var patch structures.Game
	helpers.ReadJSON(req.Body, &patch)
	defer req.Body.Close()
	fmt.Println("Received Game Patch:", patch)
	vars := mux.Vars(req)
	tid, gid := bson.ObjectIdHex(vars["tid"]), bson.ObjectIdHex(vars["gid"])
	t, _ := c.model.GetTournament(tid)
	g := t.GetGame(gid)
	g.UpdateTags(patch.Tags)
	helpers.WriteJSON(g, w)
}
