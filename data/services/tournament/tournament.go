package tournament

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
	s.store.CreateTournament(t)
	for _, g := range t.Games {
		s.store.CreateGame(g)
	}
	return t.Id, nil
}

func (s *Service) ReadTournament(id models.Id) (*models.Tournament, error) {
	t, err := s.store.ReadTournament(id)
	if err == stores.ErrNotFound {
		return nil, ErrNotFound
	}
	log.Debug("Read Tournament ", *t)
	return t, nil
}

func (s *Service) ReadTournaments(filter func(*models.Tournament) bool) []*models.Tournament {
	return s.store.ReadTournaments(filter)
}
