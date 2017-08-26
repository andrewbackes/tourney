package v2

import (
	"github.com/andrewbackes/tourney/db"
	"github.com/gorilla/mux"
	"net/http"
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
func Bind(db db.Database, r *mux.Router) {
	s := r.PathPrefix("/api/v2").Subrouter()
	h := make([][]string, 0)
	bind := func(m, p string, f func(w http.ResponseWriter, req *http.Request)) {
		s.HandleFunc(p, f).Methods(m)
		h = append(h, []string{m, p})
	}

	bind("GET", "/tournaments", func(w http.ResponseWriter, req *http.Request) {

	})
	/*
		s.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
			resp := `{`
			for k, v := range h {
				resp += `  "` +  +`":["` + +`"]`
			}
			resp += '}'
		}).Methods("GET")
	*/
}
