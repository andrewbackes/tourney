package api

import (
	"github.com/andrewbackes/tourney/data"
	"github.com/gorilla/mux"
)

// Bind sets the routes in the router.
func Bind(s data.Service, r *mux.Router) {
	sub := r.PathPrefix("/api/v2").Subrouter()
	sub.HandleFunc("/tournaments", getTournaments(s)).Methods("GET")
	sub.HandleFunc("/tournaments/{id}", getTournament(s)).Methods("GET")
	sub.HandleFunc("/tournaments", postTournament(s)).Methods("POST")
	sub.HandleFunc("/tournaments/{tid}/games/{gid}", getGame(s)).Methods("GET")
	sub.HandleFunc("/tournaments/{tid}/games/{gid}", patchGame(s)).Methods("PATCH")
	sub.HandleFunc("/tournaments/{tid}/games/{gid}/positions", postPosition(s)).Methods("POST")
}
