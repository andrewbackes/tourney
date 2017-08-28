package memdb

import (
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/stores"
	"sync"
)

func (m *MemDB) CreateGame(g *models.Game) {
	m.games.Store(g.Id, g)
	m.gameLocks.Store(g.Id, &sync.Mutex{})
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

func (m *MemDB) lock(gid models.Id) {
	lock, exists := m.gameLocks.Load(gid)
	if !exists {
		panic("missing required element in map")
	}
	lock.(*sync.Mutex).Lock()
}

func (m *MemDB) unlock(gid models.Id) {
	lock, exists := m.gameLocks.Load(gid)
	if !exists {
		panic("missing required element in map")
	}
	lock.(*sync.Mutex).Unlock()
}
