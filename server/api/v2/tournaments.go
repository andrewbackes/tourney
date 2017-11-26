package api

import (
	"fmt"
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/services"
	"github.com/andrewbackes/tourney/data/services/tournament"
	"github.com/andrewbackes/tourney/util"
	"github.com/gorilla/mux"
	//log "github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

func getTournaments(s services.Tournament) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		statusVal := req.URL.Query().Get("status")
		engineVal := req.URL.Query().Get("engine")
		filter := func(t *models.Tournament) bool {
			r := true
			if statusVal != "" && t.Status != models.Status(statusVal) {
				r = false
			}
			if engineVal != "" {
				r = false
				for i := range t.Settings.Contestants {
					if t.Settings.Contestants[i].Id() == strings.ToLower(engineVal) {
						r = true
						break
					}
				}
			}
			return r
		}
		ts := s.ReadTournaments(filter)
		collapsedTs := make([]*models.CollapsedTournament, len(ts), len(ts))
		for i, t := range ts {
			collapsedTs[i] = models.CollapseTournament(t)
		}
		util.WriteJSON(collapsedTs, w)
	}
}

func getTournament(s services.Tournament) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := models.Id(vars["id"])
		t, err := s.ReadTournament(id)
		if err == tournament.ErrNotFound {
			w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
			w.WriteHeader(http.StatusNotFound)
		} else {
			util.WriteJSON(models.CollapseTournament(t), w)
		}
	}
}

func postTournament(s services.Tournament) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		var t models.Tournament
		util.ReadJSON(req.Body, &t)
		defer req.Body.Close()
		id, _ := s.CreateTournament(&t)
		w.Write([]byte(fmt.Sprintf("{\"id\":\"%s\"}", id)))
	}
}
