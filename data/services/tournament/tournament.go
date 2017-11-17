package tournament

import (
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/stores"
	log "github.com/sirupsen/logrus"
	"time"
)

func (s *Service) CreateTournament(t *models.Tournament) (models.Id, error) {
	log.Debug("Creating Tournament ", *t)
	if len(append(t.Settings.Contestants, t.Settings.Opponents...)) < 2 {
		return "", ErrMalformed
	}
	t.Id = models.NewId()
	t.CreationDate = time.Now()
	for i := range t.Settings.Contestants {
		if t.Settings.Contestants[i].Id == "" {
			t.Settings.Contestants[i].Id = models.NewId()
		}
	}
	for i := range t.Settings.Opponents {
		if t.Settings.Opponents[i].Id == "" {
			t.Settings.Opponents[i].Id = models.NewId()
		}
	}
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
