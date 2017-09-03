package api

import (
	"github.com/andrewbackes/tourney/data"
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/service"
	"github.com/andrewbackes/tourney/util"
	"github.com/gorilla/mux"
	"net/http"
)

func getGames(s data.Service) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		id := models.Id(vars["id"])
		var filter func(*models.Game) bool
		val := req.URL.Query().Get("status")
		if val != "" {
			var status models.Status
			(&status).UnmarshalJSON([]byte(`"` + val + `"`))
			filter = func(t *models.Game) bool {
				if t.Status == status {
					return true
				}
				return false
			}
		}
		gs := s.ReadGames(id, filter)
		util.WriteJSON(gs, w)
	}
}

func getGame(s data.Service) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		tid := models.Id(vars["tid"])
		gid := models.Id(vars["gid"])
		g, err := s.ReadGame(tid, gid)
		if err == service.ErrNotFound {
			w.WriteHeader(http.StatusNotFound)
		} else {
			util.WriteJSON(g, w)
		}
	}
}

func putGame(s data.Service) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		g := &models.Game{}
		util.ReadJSON(req.Body, g)
		s.UpdateGame(g)
	}
}

/*

func postPosition(s data.Service) func(w http.ResponseWriter, req *http.Request) {
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

func patchGame(s data.Service) func(w http.ResponseWriter, req *http.Request) {
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
