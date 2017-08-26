package api

import (
	"fmt"
	"github.com/andrewbackes/tourney/data"
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/util"
	"github.com/gorilla/mux"
	"net/http"
)

func getTournaments(s data.Service) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {

	}
}

func getTournament(s data.Service) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := models.Id(vars["id"])
		t := s.ReadTournament(id)
		if t != nil {
			util.WriteJSON(t, w)
		} else {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}
}

func postTournament(s data.Service) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var t models.Tournament
		util.ReadJSON(req.Body, &t)
		defer req.Body.Close()
		id := s.CreateTournament(&t)
		w.Write([]byte(fmt.Sprintf("{\"id\":\"%s\"}", id)))
	}
}
