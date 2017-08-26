package memdb

import (
	"github.com/andrewbackes/tourney/data/models"
)

func (m *MemDB) CreateTournament(t *models.Tournament) models.Id {
	t.Id = models.NewId()
	m.tournaments.Store(t.Id, t)
	return t.Id
}

func (m *MemDB) ReadTournament(id models.Id) *models.Tournament {
	t, exists := m.tournaments.Load(id)
	if exists {
		return t.(*models.Tournament)
	}
	return nil
}

func (m *MemDB) ReadTournaments(filter func(*models.Tournament) bool) []*models.Tournament {
	result := make([]*models.Tournament, 0)
	m.tournaments.Range(func(key, value interface{}) bool {
		if filter(value.(*models.Tournament)) {
			result = append(result, value.(*models.Tournament))
		}
		return true
	})
	return nil
}

func (m *MemDB) UpdateTournamentStatus(id models.Id, status models.Status) {
	t := m.ReadTournament(id)
	t.Status = status
}

func (m *MemDB) UpdateTournamentSummary(id models.Id, summary models.Summary) {
	t := m.ReadTournament(id)
	t.Summary = summary
}
