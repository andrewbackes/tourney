package tournament

import (
	"github.com/andrewbackes/tourney/data/models"
)

func (s *Service) ReadEngines(filter func(*models.Engine) bool) []*models.Engine {
	return s.store.ReadEngines(filter)
}

func (s *Service) CreateEngine(e *models.Engine) {
	s.store.CreateEngine(e)
}
