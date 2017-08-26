package service

import (
	"github.com/andrewbackes/tourney/data"
)

type Service struct {
	store data.Store
}

func New(db data.Store) *Service {
	return &Service{
		store: db,
	}
}
