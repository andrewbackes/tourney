package service

import (
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/stores"
	log "github.com/sirupsen/logrus"
	"strings"
)

func (s *Service) ReadGame(tid, gid models.Id) (*models.Game, error) {
	g, err := s.store.ReadGame(tid, gid)
	log.Debug("Read Game ", *g)
	if err == stores.ErrNotFound {
		return nil, ErrNotFound
	}
	return g, nil
}

func (s *Service) AddPosition(tid, gid models.Id, p models.Position) error {
	a := strings.Split(p.FEN, " ")
	if len(a) < 6 {
		return ErrMalformed
	}
	s.store.AddPosition(tid, gid, p)
	return nil
}

func (s *Service) CompleteGame(tid, gid models.Id) {
	s.store.UpdateStatus(tid, gid, models.Complete)
}

func (s *Service) AssignGame(tid, gid models.Id) {
	s.store.UpdateStatus(tid, gid, models.Running)
}
