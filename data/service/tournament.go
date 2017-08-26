package service

import (
	"github.com/andrewbackes/tourney/data/models"
)

func (s *Service) CreateTournament(t *models.Tournament) (models.Id, error) {
	id := models.NewId()
	if len(t.Settings.Engines) < 2 {
		return "", ErrMalformed
	}
	t.Games = models.NewGameList(id, t.Settings)
	s.store.CreateTournament(t)
	return id, nil
}

func (s *Service) ReadTournament(id models.Id) (*models.Tournament, error) {
	return s.store.ReadTournament(id), nil
}
