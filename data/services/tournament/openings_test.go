package tournament

import (
	"fmt"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/tourney/data/models"
	"testing"
)

func TestSetOpenings(t *testing.T) {
	trnmt := newTournament()
	trnmt.Games = models.NewGameList("tId", trnmt.Settings) // TODO: mock this
	setGameOpenings(trnmt)
	fmt.Println(trnmt.Games)
}

func newTournament() *models.Tournament {
	return &models.Tournament{
		Settings: models.Settings{
			TestSeats: 1,
			Carousel:  false,
			Rounds:    10,
			Engines: []models.Engine{
				models.Engine{
					Name: "tester1",
				},
				models.Engine{
					Name: "tester2",
				},
			},
			TimeControl: game.TimeControl{
				Moves:     40,
				Time:      1000000000,
				Repeating: true,
			},
			Opening: models.Opening{
				Depth:     8,
				Randomize: true,
				Book: models.Book{
					FilePath: "/Users/Andrew/tourney_books/2700draw.bin",
					MaxDepth: 14,
				},
			},
		},
	}
}
