package api

import (
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/services"
	"github.com/andrewbackes/tourney/data/services/tournament"
	"github.com/andrewbackes/tourney/util"
	"github.com/gorilla/mux"
	"net/http"
)

func getGames(s services.Tournament) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := models.Id(vars["id"])
		var filter func(*models.Game) bool
		val := req.URL.Query().Get("status")
		if val != "" {
			//var status models.Status
			//(&status).UnmarshalJSON([]byte(`"` + val + `"`))
			filter = func(t *models.Game) bool {
				if t.Status == models.Status(val) {
					return true
				}
				return false
			}
		}
		gs := s.ReadGames(id, filter)
		collapsedGs := make([]*models.CollapsedGame, len(gs), len(gs))
		for i, g := range gs {
			collapsedGs[i] = models.CollapseGame(g)
		}
		util.WriteJSON(collapsedGs, w)
	}
}

func getGame(s services.Tournament) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tid := models.Id(vars["tid"])
		gid := models.Id(vars["gid"])
		g, err := s.ReadGame(tid, gid)
		if err == tournament.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			util.WriteJSON(g, w)
		}
	}
}

func putGame(s services.Tournament) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		g := &models.Game{}
		util.ReadJSON(req.Body, g)
		s.UpdateGame(g)
	}
}

/*

func postPosition(s services.Tournament) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tid := models.Id(vars["tid"])
		gid := models.Id(vars["gid"])
		var p models.Position
		util.ReadJSON(req.Body, &p)
		defer req.Body.Close()
		s.AddPosition(tid, gid, p)
		w.Write([]byte(fmt.Sprintf("{\"status\":\"%s\"}", "success")))
	}
}

func patchGame(s services.Tournament) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tid := models.Id(vars["tid"])
		gid := models.Id(vars["gid"])
		patch := struct {
			Status *models.Status
		}{}
		util.ReadJSON(req.Body, &patch)
		if patch.Status != nil {
			switch *patch.Status {
			case models.Complete:
				s.CompleteGame(tid, gid)
				w.Write([]byte(fmt.Sprintf("{\"status\":\"%s\"}", "success")))
			case models.Running:
				s.AssignGame(tid, gid)
				w.Write([]byte(fmt.Sprintf("{\"status\":\"%s\"}", "success")))
			default:
				w.WriteHeader(http.StatusBadRequest)
			}
		}
	}
}

*/
