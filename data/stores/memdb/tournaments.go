package memdb

import (
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/data/stores"
)

func (m *MemDB) CreateTournament(t *models.Tournament) {
	m.tournaments.Store(t.Id, t)
}

func (m *MemDB) ReadTournament(id models.Id) (*models.Tournament, error) {
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

func (m *MemDB) UpdateTournamentStatus(id models.Id, status models.Status) {
	t, _ := m.ReadTournament(id)
	t.Status = status
}

func (m *MemDB) UpdateTournamentSummary(id models.Id, summary models.Summary) {
	t, _ := m.ReadTournament(id)
	t.Summary = summary
}
