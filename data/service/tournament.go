package service

import (
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/stores"
	log "github.com/sirupsen/logrus"
)

func (s *Service) CreateTournament(t *models.Tournament) (models.Id, error) {
	log.Debug("Creating Tournament ", *t)
	if len(t.Settings.Engines) < 2 {
		return "", ErrMalformed
	}
	t.Id = models.NewId()
	t.Games = models.NewGameList(t.Id, t.Settings)
	id := s.store.CreateTournament(t)
	return id, nil
}

func (s *Service) ReadTournament(id models.Id) (*models.Tournament, error) {
	t, err := s.store.ReadTournament(id)
	log.Debug("Read Tournament ", *t)
	if err == stores.ErrNotFound {
		return nil, ErrNotFound
	}
	return t, nil
}

func (s *Service) ReadTournaments(filter func(*models.Tournament) bool) []*models.Tournament {
	return s.store.ReadTournaments(filter)
}
