package data

import (
	"github.com/andrewbackes/tourney/data/models"
)

// Service contains logic for handling changes in data models.
type Service interface {
	CreateTournament(*models.Tournament) (models.Id, error)
	ReadTournament(id models.Id) (*models.Tournament, error)
	ReadTournaments(filter func(*models.Tournament) bool) []*models.Tournament

	ReadGame(tid, gid models.Id) (*models.Game, error)
	AddPosition(tid, gid models.Id, p models.Position) error

	CompleteGame(tid, gid models.Id)
	AssignGame(tid, gid models.Id)
}
