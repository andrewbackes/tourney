package data

import (
	"github.com/andrewbackes/tourney/data/models"
)

// Store persists data vai CRUD.
type Store interface {
	CreateTournament(*models.Tournament) models.Id
	CreateGame(*models.Game) models.Id

	ReadTournament(id models.Id) *models.Tournament
	ReadTournaments(filter func(*models.Tournament) bool) []*models.Tournament
	ReadGame(id models.Id) *models.Game
	ReadGames(filter func(*models.Game) bool) []*models.Game

	UpdateTournamentSummary(id models.Id, summary models.Summary)
	UpdateTournamentStatus(id models.Id, status models.Status)
	UpdateGameTags(id models.Id, tags map[string]string)
	UpdateGameStatus(id models.Id, status models.Status)
	UpdateGamePosition(id models.Id)
}
