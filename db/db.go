package db

import (
	"github.com/andrewbackes/tourney/models"
)

// Database is what is needed to CRUD the models.
type Database interface {
	createTournament(*models.Tournament) models.Id
	createGame(*models.Game) models.Id

	readTournament(id models.Id) *models.Tournament
	readTournaments(filter func(*models.Tournament) bool) []*models.Tournament
	readGame(id models.Id) *models.Game
	readGames(filter func(*models.Game) bool) []*models.Game

	updateTournamentSummary(id models.Id, summary models.Summary)
	updateTournamentStatus(id models.Id, status models.Status)
	updateGameTags(id models.Id, tags map[string]string)
	updateGameStatus(id models.Id, status models.Status)
	updateGamePosition(id models.Id)
}
