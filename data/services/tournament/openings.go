package tournament

import (
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
// Assumes that games come in matchup pairs.
// TODO(andrewbackes): increment
// TODO(andrewbackes): panic handling
func setGameOpenings(t *models.Tournament) {
	if t.Settings.TimeControl.Moves < t.Settings.Opening.Depth {
		panic("book moves >= moves per time control")
	}
	if t.Settings.Opening.Depth > t.Settings.Opening.Book.MaxDepth {
		panic("book depth is less than what is specified in the tournament")
	}
	b := openBook(t.Settings.Opening.Book.FilePath)
	o := openingPos(t.Settings.TimeControl)
	pl := []*models.Position{posToPos(o)}
	index := 0
	complete := false
	deepen(t.Settings.Opening.Depth, o, pl, b, complete, func(l []*models.Position) {
		if !complete {
			for j := 0; j < 2; j++ {
				t.Games[index].Positions = posPtrsToVals(l)
				index++
				if index >= len(t.Games)-1 {
					complete = true
					break
				}
			}
		}
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

func deepen(d int, orig *position.Position, pl []*models.Position, b *book.Book, complete bool, callback func([]*models.Position)) {
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
			deepen(d-1, next, nextSlice, b, complete, callback)
			if complete {
				return
			}
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

func posPtrsToVals(ptrs []*models.Position) []models.Position {
	n := make([]models.Position, len(ptrs))
	for i, p := range ptrs {
		n[i] = models.Position{
			FEN:      p.FEN,
			LastMove: p.LastMove,
			MovesLeft: map[piece.Color]int{
				piece.White: p.MovesLeft[piece.White],
				piece.Black: p.MovesLeft[piece.Black],
			},
			Clocks: map[piece.Color]time.Duration{
				piece.White: p.Clocks[piece.White],
				piece.Black: p.Clocks[piece.Black],
			},
		}
	}
	return n
}
