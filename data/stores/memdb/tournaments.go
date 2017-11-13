package memdb

import (
	"encoding/json"
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/stores"
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"sync"
)

func (m *MemDB) CreateTournament(t *models.Tournament) {
	m.tournaments.Store(t.Id, t)
	m.locks.Store(t.Id, &sync.Mutex{})
	m.lock(t.Id)
	defer m.unlock(t.Id)
	m.persistTournament(t)
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
	m.persistTournament(t)
}

func (m *MemDB) persistTournament(t *models.Tournament) {
	if !m.persisted() {
		return
	}
	val := *t
	// clear out the games so that they don't get saved in the json
	val.Games = nil
	tournamentDir := filepath.Join(m.backupDir, "tournaments")
	err := os.MkdirAll(tournamentDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	tournamentJSON := filepath.Join(tournamentDir, t.Id.String()+".json")
	f, err := os.Create(tournamentJSON)
	if err != nil {
		panic(err)
	}
	log.Info("Persisting tournament", tournamentJSON)
	err = json.NewEncoder(f).Encode(t)
	if err != nil {
		panic(err)
	}
}
