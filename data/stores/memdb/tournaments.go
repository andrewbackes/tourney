package memdb

import (
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/stores"
	log "github.com/sirupsen/logrus"
	"sync"
)

func (m *MemDB) CreateTournament(t *models.Tournament) {
	m.tournaments.Store(t.Id, t)
	m.locks.Store(t.Id, &sync.Mutex{})
}

func (m *MemDB) ReadTournament(id models.Id) (*models.Tournament, error) {
	m.lock(id)
	defer m.unlock(id)
	return m.readTournament(id)
}

func (m *MemDB) readTournament(id models.Id) (*models.Tournament, error) {
	t, exists := m.tournaments.Load(id)
	if !exists {
		return nil, stores.ErrNotFound
	}
	return t.(*models.Tournament), nil
}

func (m *MemDB) ReadTournaments(filter func(*models.Tournament) bool) []*models.Tournament {
	result := make([]*models.Tournament, 0)
	m.tournaments.Range(func(id, tournament interface{}) bool {
		if filter == nil || filter(tournament.(*models.Tournament)) {
			result = append(result, tournament.(*models.Tournament))
		}
		return true
	})
	return result
}

func (m *MemDB) UpdateTournament(t *models.Tournament) {
	m.lock(t.Id)
	defer m.unlock(t.Id)
	old, err := m.readTournament(t.Id)
	if err != nil {
		log.Error("Could not read tournament ", t.Id, ": ", err)
	}
	*old = *t
}
