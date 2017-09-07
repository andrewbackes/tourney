package models

import (
	"fmt"
	"github.com/andrewbackes/chess/game"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSetOpenings(t *testing.T) {
	trnmt := newTournament()
	trnmt.Games = NewGameList("tId", trnmt.Settings) // TODO: mock this
	fmt.Println(trnmt.Games)
	assert.Equal(t, 20, len(trnmt.Games))
}

func newTournament() *Tournament {
	return &Tournament{
		Settings: Settings{
			TestSeats: 1,
			Carousel:  false,
			Rounds:    10,
			Engines: []Engine{
				Engine{
					Name: "tester1",
				},
				Engine{
					Name: "tester2",
				},
			},
			TimeControl: game.TimeControl{
				Moves:     40,
				Time:      1000000000,
				Repeating: true,
			},
			Opening: Opening{
				Depth:     8,
				Randomize: true,
				Book: Book{
					FilePath: "/Users/Andrew/tourney_books/2700draw.bin",
					MaxDepth: 14,
				},
			},
		},
	}
}
