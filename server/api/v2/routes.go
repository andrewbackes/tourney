package api

import (
	"github.com/andrewbackes/tourney/data/services"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
)

// Bind sets the routes in the router.
func Bind(s services.Tournament, r *mux.Router) {
	sub := r.PathPrefix("/api/v2").Subrouter()
	register := func(method, path string, f func(services.Tournament) func(http.ResponseWriter, *http.Request)) {
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
	register("PUT", "/tournaments/{tid}/games/{gid}", putGame)
	//register("GET", "/engines}" getEngines)
	//register("GET", "/engines/{name}/{version}/{os}" getEngine)
	//register("PUT", "/engines/{name}/{version}/{os}" putEngine)
	register("POST", "/engineFiles/{name}/{version}/{os}", postEngineFile)
	register("GET", "/engineFiles/{name}/{version}/{os}", getEngineFile)
	register("POST", "/engines", postEngine)
	register("GET", "/engines", getEngines)
}
