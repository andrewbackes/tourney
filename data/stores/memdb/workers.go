package memdb

import (
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/stores"
)

func (m *MemDB) CreateWorker(w models.Worker) {
	m.workers.Store(w.Id, w)
}

func (m *MemDB) ReadWorker(id models.Id) (models.Worker, error) {
	w, exists := m.workers.Load(id)
	if exists {
		return w.(models.Worker), nil
	}
	return models.Worker{}, stores.ErrNotFound
}

func (m *MemDB) DeleteWorker(id models.Id) {
	m.workers.Delete(id)
}
