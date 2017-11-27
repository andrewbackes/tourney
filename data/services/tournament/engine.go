package tournament

import (
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/stores"
)

func (s *Service) ReadEngines(filter func(*models.Engine) bool) []*models.Engine {
	return s.store.ReadEngines(filter)
}

func (s *Service) CreateEngine(e *models.Engine) {
	s.store.CreateEngine(e)
}

func (s *Service) ReadEngine(id string) (*models.Engine, error) {
	return s.store.ReadEngine(id)
}

func (s *Service) EngineExists(id string) bool {
	_, err := s.store.ReadEngine(id)
	if err == stores.ErrNotFound {
		return false
	}
	return true
}
