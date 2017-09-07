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

func (m *MemDB) CreateGame(g *models.Game) {
	m.games.Store(g.Id, g)
	m.locks.Store(g.Id, &sync.Mutex{})
	m.lock(g.Id)
	defer m.unlock(g.Id)
	m.persistGame(g)
}

func (m *MemDB) UpdateGame(g *models.Game) {
	old, err := m.ReadGame(g.TournamentId, g.Id)
	if err != nil {
		log.Error("Could not read game from data store: ", err)
	}
	m.lock(g.Id)
	defer m.unlock(g.Id)
	*old = *g
	m.persistGame(g)
}

func (m *MemDB) ReadGame(tid, gid models.Id) (*models.Game, error) {
	m.lock(gid)
	defer m.unlock(gid)
	return m.readGame(tid, gid)
}

func (m *MemDB) readGame(tid, gid models.Id) (*models.Game, error) {
	g, exists := m.games.Load(gid)
	if exists {
		return g.(*models.Game), nil
	}
	return nil, stores.ErrNotFound
}

func (m *MemDB) ReadGames(tid models.Id, filter func(*models.Game) bool) []*models.Game {
	t, err := m.ReadTournament(tid)
	result := make([]*models.Game, 0)
	if err != nil {
		log.Error("Could not read tournament ", tid, " : ", err)
		return result
	}
	if filter == nil {
		return t.Games
	}
	for _, g := range t.Games {
		if filter(g) {
			result = append(result, g)
		}
	}
	return result
}

func (m *MemDB) persistGame(g *models.Game) {
	if !m.persisted() || g.Status == models.Running {
		return
	}
	tournamentDir := filepath.Join(m.backupDir, "tournaments", g.TournamentId.String())
	err := os.MkdirAll(tournamentDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	gameJSON := filepath.Join(tournamentDir, g.Id.String()+".json")
	f, err := os.Create(gameJSON)
	if err != nil {
		panic(err)
	}
	log.Info("Persisting game ", gameJSON)
	err = json.NewEncoder(f).Encode(g)
	if err != nil {
		panic(err)
	}
}

/*
func (m *MemDB) AddPosition(tid, gid models.Id, p models.Position) error {
	m.lock(gid)
	defer m.unlock(gid)
	ply := (p.MoveNumber() - 1) * 2
	if p.ActiveColor() == piece.Black {
		ply++
	}
	g, err := m.readGame(tid, gid)
	if err != nil {
		return err
	}
	for len(g.Positions) < ply+1 {
		g.Positions = append(g.Positions, models.Position{})
	}
	g.Positions[ply] = p
	return nil
}

func (m *MemDB) UpdateStatus(tid, gid models.Id, s models.Status) {
	m.lock(gid)
	defer m.unlock(gid)
	g, _ := m.readGame(tid, gid)
	g.Status = s
}
*/
