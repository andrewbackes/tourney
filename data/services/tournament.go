package services

import (
	"github.com/andrewbackes/tourney/data/models"
)

// Tournament contains logic for handling changes in data models.
type Tournament interface {
	CreateTournament(*models.Tournament) (models.Id, error)
	ReadTournament(id models.Id) (*models.Tournament, error)
	ReadTournaments(filter func(*models.Tournament) bool) []*models.Tournament

	ReadGame(tid, gid models.Id) (*models.Game, error)
	ReadGames(tid models.Id, filter func(*models.Game) bool) []*models.Game
	UpdateGame(*models.Game) error
}
