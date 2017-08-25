package server

import (
	"github.com/andrewbackes/tourney/db"
)

type Server struct{}

func NewServer(port string, db db.Database) *Server {
	return &Server{}
}

func (s *Server) Start() {}
