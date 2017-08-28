package api

import (
	"github.com/andrewbackes/tourney/data"
	"github.com/gorilla/mux"
)

/*
func temp(w http.ResponseWriter, req *http.Request) {}

func newRouter(c *controller) *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v2").Subrouter()
	s.HandleFunc("/tournaments", c.getTournaments).Methods("GET")
	s.HandleFunc("/tournaments", c.postTournament).Methods("POST")
	s.HandleFunc("/tournaments/{tid}", c.getTournament).Methods("GET")
	s.HandleFunc("/tournaments/{tid}/games", c.getGames).Methods("GET")
	s.HandleFunc("/tournaments/{tid}/games/{gid}", c.getGame).Methods("GET")
	s.HandleFunc("/tournaments/{tid}/games/{gid}", c.patchGame).Methods("PATCH")
	s.HandleFunc("/tournaments/{tid}/games/{gid}/positions", c.getPositions).Methods("GET")
	s.HandleFunc("/tournaments/{tid}/games/{gid}/positions", c.postPosition).Methods("POST")
	s.HandleFunc("/tournaments/{tid}/games/{gid}/positions/{pid}", c.getPosition).Methods("GET")
	s.HandleFunc("/tournaments/{tid}/gameQueue/next", c.getTournamentsNextGame).Methods("GET")
	s.HandleFunc("/workers", c.postWorker).Methods("POST")
	return r
}
*/

// Bind sets the routes in the router.
func Bind(s data.Service, r *mux.Router) {
	sub := r.PathPrefix("/api/v2").Subrouter()
	sub.HandleFunc("/tournaments", getTournaments(s)).Methods("GET")
	sub.HandleFunc("/tournaments/{id}", getTournament(s)).Methods("GET")
	sub.HandleFunc("/tournaments", postTournament(s)).Methods("POST")
}
