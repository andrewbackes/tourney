package data

import (
	"github.com/andrewbackes/tourney/data/models"
)

// Service contains logic for handling changes in data models.
type Service interface {
	CreateTournament(*models.Tournament) (models.Id, error)
	ReadTournament(id models.Id) (*models.Tournament, error)
	ReadTournaments(filter func(*models.Tournament) bool) []*models.Tournament
}
