package tournament

import (
	"fmt"
	"github.com/andrewbackes/chess/book"
	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"github.com/andrewbackes/tourney/data/models"
	"os"
	"time"
)

// setGameOpenings will set the state of all games in the tournament. It chooses the state based on the opening book selected.
// TODO(andrewbackes): increment
// TODO(andrewbackes): panic handling
func setGameOpenings(t *models.Tournament) {
	if t.Settings.TimeControl.Moves < t.Settings.Opening.Depth {
		panic("book moves >= moves per time control")
	}
	if t.Settings.Opening.Depth > t.Settings.Opening.Book.MaxDepth {
		panic("book depth is less than what is specified in the tournament")
	}
	b := openBook("/Users/Andrew/tourney_books/2700draw.bin")
	o := openingPos(t.Settings.TimeControl)
	pl := []*models.Position{posToPos(o)}
	deepen(t.Settings.Opening.Depth, o, pl, b, func(pl []*models.Position) {
		fmt.Println(pl)
	})
}

func openBook(path string) *book.Book {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	b, err := book.Read(f)
	if err != nil {
		panic(err)
	}
	return b
}

func openingPos(tc game.TimeControl) *position.Position {
	o, err := fen.Decode("rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1")
	if err != nil {
		panic(err)
	}
	o.Clocks = map[piece.Color]time.Duration{
		piece.White: tc.Time,
		piece.Black: tc.Time,
	}
	o.MovesLeft = map[piece.Color]int{
		piece.White: tc.Moves,
		piece.Black: tc.Moves,
	}
	return o
}

func deepen(d int, orig *position.Position, pl []*models.Position, b *book.Book, callback func([]*models.Position)) {
	if d == 0 {
		callback(pl)
		return
	}
	if entries, exists := b.Positions[orig.Polyglot()]; exists {
		for _, entry := range entries {
			next := orig.MakeMove(entry.Move)
			nextSlice := make([]*models.Position, len(pl))
			copy(nextSlice, pl)
			nextSlice = append(nextSlice, posToPos(next))
			deepen(d-1, next, nextSlice, b, callback)
		}
	}
}

func posToPos(orig *position.Position) *models.Position {
	f, err := fen.Encode(orig)
	if err != nil {
		panic(err)
	}
	return &models.Position{
		FEN: f,
		MovesLeft: map[piece.Color]int{
			piece.White: orig.MovesLeft[piece.White],
			piece.Black: orig.MovesLeft[piece.Black],
		},
		Clocks: map[piece.Color]time.Duration{
			piece.White: orig.Clocks[piece.White],
			piece.Black: orig.Clocks[piece.Black],
		},
		LastMove: orig.LastMove,
	}
}
