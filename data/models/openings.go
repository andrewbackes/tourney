package models

import (
	"github.com/andrewbackes/chess/book"
	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position"
	"os"
	"time"
)

// setGameOpenings will set the state of all games in the tournament. It chooses the state based on the opening book selected.
// Assumes that games come in matchup pairs.
func setGameOpenings(list []*Game, settings Settings) {
	if settings.TimeControl.Moves < settings.Opening.Depth {
		panic("book moves >= moves per time control")
	}
	if settings.Opening.Depth > settings.Opening.Book.MaxDepth {
		panic("book depth is less than what is specified in the tournament")
	}
	b := openBook(settings.Opening.Book.FilePath)
	o := openingPos(settings.TimeControl)
	pl := []*Position{posToPos(o)}
	index := 0
	complete := false
	deepen(settings.Opening.Depth, o, pl, b, complete, func(l []*Position) {
		if !complete {
			for j := 0; j < 2; j++ {
				list[index].Positions = posPtrsToVals(l)
				index++
				if index >= len(list) {
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

func deepen(d int, orig *position.Position, pl []*Position, b *book.Book, complete bool, callback func([]*Position)) {
	if d == 0 {
		callback(pl)
		return
	}
	if entries, exists := b.Positions[orig.Polyglot()]; exists {
		for _, entry := range entries {
			next := orig.MakeMove(entry.Move)
			nextSlice := make([]*Position, len(pl))
			copy(nextSlice, pl)
			nextSlice = append(nextSlice, posToPos(next))
			deepen(d-1, next, nextSlice, b, complete, callback)
			if complete {
				return
			}
		}
	}
}

func posToPos(orig *position.Position) *Position {
	f, err := fen.Encode(orig)
	if err != nil {
		panic(err)
	}
	return &Position{
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

func posPtrsToVals(ptrs []*Position) []Position {
	n := make([]Position, len(ptrs))
	for i, p := range ptrs {
		n[i] = Position{
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
