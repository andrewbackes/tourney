package server

import (
	"github.com/andrewbackes/tourney/data"
	"github.com/andrewbackes/tourney/server/api/v2"
	"github.com/andrewbackes/tourney/server/ui"
	"github.com/gorilla/mux"
)

type Server struct {
	router *mux.Router
	port   string
}

func New(port string, s data.Service) *Server {
	r := mux.NewRouter()
	api.Bind(s, r)
	ui.Bind(r)
	return &Server{
		router: r,
		port:   port,
	}
}

func (s *Server) Start() {}
