package api

import (
	"github.com/andrewbackes/tourney/data"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Bind sets the routes in the router.
func Bind(s data.Service, r *mux.Router) {
	sub := r.PathPrefix("/api/v2").Subrouter()
	register := func(method, path string, f func(data.Service) func(http.ResponseWriter, *http.Request)) {
		wrapper := func() func(http.ResponseWriter, *http.Request) {
			return func(w http.ResponseWriter, req *http.Request) {
				log.Debug("[", req.RemoteAddr, "] ", req.Method, " ", req.RequestURI)
				f(s)(w, req)
			}
		}
		sub.HandleFunc(path, wrapper()).Methods(method)
	}
	register("GET", "/tournaments", getTournaments)
	register("GET", "/tournaments/{id}", getTournament)
	register("POST", "/tournaments", postTournament)
	register("GET", "/tournaments/{id}/games", getGames)
	register("GET", "/tournaments/{tid}/games/{gid}", getGame)
	register("PATCH", "/tournaments/{tid}/games/{gid}", patchGame)
	register("POST", "/tournaments/{tid}/games/{gid}/positions", postPosition)
}
