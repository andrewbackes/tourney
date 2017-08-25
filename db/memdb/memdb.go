package memdb

import (
	"github.com/andrewbackes/tourney/models"
)

type MemDB struct{}

func NewMemDB() *MemDB {
	return &MemDB{}
}

func (m *MemDB) createTournament(*models.Tournament) models.Id {
	return ""
}

func (m *MemDB) createGame(*models.Game) models.Id {
	return ""
}

func (m *MemDB) readTournament(id models.Id) *models.Tournament {
	return nil
}

func (m *MemDB) readTournaments(filter func(*models.Tournament) bool) []*models.Tournament {
	return nil
}

func (m *MemDB) readGame(id models.Id) *models.Game {
	return nil
}

func (m *MemDB) readGames(filter func(*models.Game) bool) []*models.Game {
	return nil
}

func (m *MemDB) updateTournamentSummary(id models.Id, summary models.Summary) {}
func (m *MemDB) updateTournamentStatus(id models.Id, status models.Status)    {}
func (m *MemDB) updateGameTags(id models.Id, tags map[string]string)          {}
func (m *MemDB) updateGameStatus(id models.Id, status models.Status)          {}
func (m *MemDB) updateGamePosition(id models.Id)                              {}
