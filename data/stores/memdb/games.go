package memdb

import (
	"github.com/andrewbackes/tourney/data/models"
)

func (m *MemDB) CreateGame(g *models.Game) models.Id {
	g.Id = models.NewId()
	m.games.Store(g.Id, g)
	return ""
}

func (m *MemDB) ReadGame(id models.Id) *models.Game {
	g, exists := m.games.Load(id)
	if exists {
		return g.(*models.Game)
	}
	return nil
}

func (m *MemDB) ReadGames(filter func(*models.Game) bool) []*models.Game {
	result := make([]*models.Game, 0)
	m.games.Range(func(key, value interface{}) bool {
		if filter(value.(*models.Game)) {
			result = append(result, value.(*models.Game))
		}
		return true
	})
	return nil
}

func (m *MemDB) UpdateGameStatus(id models.Id, status models.Status) {
	g := m.ReadGame(id)
	g.Status = status
}

func (m *MemDB) UpdateGameTags(id models.Id, tags map[string]string) {
	// TODO: mutex lock
	/*
		g := m.ReadGame(id)
		for k, v := range tags {
			g.Tags[k] = v
		}
	*/
}

func (m *MemDB) UpdateGamePosition(id models.Id) {
	// TODO
}
