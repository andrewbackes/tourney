package api

import (
	"github.com/gorilla/mux"
	"net/http"
)

func temp(w http.ResponseWriter, req *http.Request) {}

func router(c *controller) *mux.Router {
	r := mux.NewRouter()
	s := r.PathPrefix("/api/v2").Subrouter()
	s.HandleFunc("/tournaments", c.getTournaments).Methods("GET")
	s.HandleFunc("/tournaments", c.postTournament).Methods("POST")
	s.HandleFunc("/tournaments/{id}", c.getTournament).Methods("GET")
	s.HandleFunc("/workers", c.postWorker).Methods("POST")
	return r
}
