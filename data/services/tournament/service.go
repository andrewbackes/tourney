package tournament

import (
	"github.com/andrewbackes/tourney/data"
)

type Service struct {
	store data.Store
}

func NewService(db data.Store) *Service {
	return &Service{
		store: db,
	}
}
