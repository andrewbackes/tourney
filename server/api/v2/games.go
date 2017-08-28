package api

import (
	"fmt"
	"github.com/andrewbackes/tourney/data"
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/service"
	"github.com/andrewbackes/tourney/util"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

func getGame(s data.Service) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Debug("Request recieved: ", *req)
		vars := mux.Vars(req)
		tid := models.Id(vars["tid"])
		gid := models.Id(vars["gid"])
		g, err := s.ReadGame(tid, gid)
		if err == service.ErrNotFound {
			w.Write([]byte(fmt.Sprintf(`{"error":"%s"}`, err)))
			w.WriteHeader(http.StatusNotFound)
		} else {
			util.WriteJSON(g, w)
		}
	}
}

func postPosition(s data.Service) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Debug("Request recieved: ", *req)
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
		log.Debug("Request recieved: ", *req)
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
			case models.Running:
				s.AssignGame(tid, gid)
			}
		}
	}
}
