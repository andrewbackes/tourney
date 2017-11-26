package tournament

import (
	"github.com/andrewbackes/tourney/data/models"
)

func (s *Service) ReadEngines() []*models.Engine {
	return s.store.ReadEngines()
}

func (s *Service) CreateEngine(e *models.Engine) {
	s.store.CreateEngine(e)
}
