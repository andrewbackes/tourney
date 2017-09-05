package tournament

import (
	"github.com/andrewbackes/tourney/data/models"
)

func (s *Service) CreateWorker(w models.Worker) models.Id {
	w.Id = models.NewId()
	return w.Id
}

func (s *Service) DeleteWorker(id models.Id) {

}
