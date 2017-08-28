package memdb

import (
	"github.com/andrewbackes/tourney/data/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAddPosition(t *testing.T) {
	m := NewMemDB()
	tid, gid := models.Id("0"), models.Id("1")
	m.CreateGame(&models.Game{
		Id:           gid,
		TournamentId: tid,
		Positions:    make([]models.Position, 0),
	})
	opening := "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1"
	p := models.Position{
		FEN: opening,
	}
	m.AddPosition(tid, gid, p)
	g, err := m.ReadGame(tid, gid)
	assert.Equal(t, nil, err)
	assert.Equal(t, opening, g.Positions[0].FEN)
}
