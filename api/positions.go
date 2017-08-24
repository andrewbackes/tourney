package api

import (
	"github.com/andrewbackes/tourney/helpers"
	"github.com/andrewbackes/tourney/model/structures"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
)

func (c *controller) getPositions(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tid := bson.ObjectIdHex(vars["tid"])
	gid := bson.ObjectIdHex(vars["gid"])
	t, _ := c.model.GetTournament(tid)
	g := t.GetGame(gid)
	p := g.GetPositions
	helpers.WriteJSON(p, w)
}

func (c *controller) getPosition(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	tid := bson.ObjectIdHex(vars["tid"])
	gid := bson.ObjectIdHex(vars["gid"])
	pid, _ := strconv.Atoi(vars["pid"])
	t, _ := c.model.GetTournament(tid)
	g := t.GetGame(gid)
	p := g.GetPosition(pid)
	helpers.WriteJSON(p, w)
}

func (c *controller) postPosition(w http.ResponseWriter, req *http.Request) {
	var p structures.Position
	helpers.ReadJSON(req.Body, &p)
	defer req.Body.Close()
	vars := mux.Vars(req)
	tid := bson.ObjectIdHex(vars["tid"])
	gid := bson.ObjectIdHex(vars["gid"])
	t, _ := c.model.GetTournament(tid)
	g := t.GetGame(gid)
	g.AddPosition(&p)
	w.Write([]byte("{\"status\":\"success\"}"))
}
