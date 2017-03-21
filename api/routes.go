package api

import (
	"github.com/gorilla/mux"
)

func router(c *controller) *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v2").Subrouter()
	s.HandleFunc("/tournaments", c.getTournaments).Methods("GET")
	s.HandleFunc("/tournaments", c.postTournaments).Methods("POST")
	s.HandleFunc("/tournaments/{id}", c.postTournaments).Methods("GET")
	return r
}
