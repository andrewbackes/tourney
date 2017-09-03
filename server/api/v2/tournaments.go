package api

import (
	"fmt"
	"github.com/andrewbackes/tourney/data"
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/service"
	"github.com/andrewbackes/tourney/util"
	"github.com/gorilla/mux"
	"net/http"
)

func getTournaments(s data.Service) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var filter func(*models.Tournament) bool
		val := req.URL.Query().Get("status")
		if val != "" {
			var status models.Status
			(&status).UnmarshalJSON([]byte(`"` + val + `"`))
			filter = func(t *models.Tournament) bool {
				if t.Status == status {
					return true
				}
				return false
			}
		}
		ts := s.ReadTournaments(filter)
		util.WriteJSON(ts, w)
	}
}

func getTournament(s data.Service) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := models.Id(vars["id"])
		t, err := s.ReadTournament(id)
		if err == service.ErrNotFound {
			w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
			w.WriteHeader(http.StatusNotFound)
		} else {
			util.WriteJSON(t, w)
		}
	}
}

func postTournament(s data.Service) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var t models.Tournament
		util.ReadJSON(req.Body, &t)
		defer req.Body.Close()
		id, _ := s.CreateTournament(&t)
		w.Write([]byte(fmt.Sprintf("{\"id\":\"%s\"}", id)))
	}
}
