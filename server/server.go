package server

import (
	"github.com/andrewbackes/tourney/data/services"
	"github.com/andrewbackes/tourney/server/api/v2"
	"github.com/andrewbackes/tourney/server/ui"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"net/http"
)

type Server struct {
	router *mux.Router
	port   string
}

func New(port string, s services.Tournament) *Server {
	r := mux.NewRouter()
	api.Bind(s, r)
	ui.Bind(r)
	return &Server{
		router: r,
		port:   port,
	}
}

func (s *Server) Start() {
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})
	http.ListenAndServe(s.port, handlers.CORS(originsOk, headersOk, methodsOk)(s.router))
}
