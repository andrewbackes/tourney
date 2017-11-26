package data

import (
	"github.com/andrewbackes/tourney/data/models"
)

// Store persists data vai CRUD.
type Store interface {
	CreateTournament(*models.Tournament)
	ReadTournament(id models.Id) (*models.Tournament, error)
	ReadTournaments(filter func(*models.Tournament) bool) []*models.Tournament
	UpdateTournament(*models.Tournament)

	CreateGame(*models.Game)
	ReadGame(tid, gid models.Id) (*models.Game, error)
	ReadGames(tid models.Id, filter func(*models.Game) bool) []*models.Game
	UpdateGame(*models.Game)

	ReadEngines(filter func(*models.Engine) bool) []*models.Engine
	CreateEngine(*models.Engine)
}
