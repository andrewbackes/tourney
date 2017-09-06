package tournament

import (
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/tourney/data/models"
	"testing"
	"time"
)

func TestOpenings(t *testing.T) {
	goDepth(game.TimeControl{Moves: 40, Time: 10 * time.Second}, models.Book{}, 2)
}
