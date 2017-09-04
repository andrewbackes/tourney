package models

import (
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"strconv"
	"strings"
	"time"
)

type Position struct {
	FEN       string                        `json:"fen"`
	LastMove  move.Move                     `json:"lastMove"`
	MovesLeft map[piece.Color]int           `json:"movesLeft"`
	Clocks    map[piece.Color]time.Duration `json:"clock"`
}

func StartPosition(c game.TimeControl) Position {
	return Position{
		FEN:       "rnbqkbnr/pppppppp/8/8/8/8/PPPPPPPP/RNBQKBNR w KQkq - 0 1",
		LastMove:  move.Null,
		MovesLeft: map[piece.Color]int{piece.White: c.Moves, piece.Black: c.Moves},
		Clocks:    map[piece.Color]time.Duration{piece.White: c.Time, piece.Black: c.Time},
	}
}

func (p *Position) MoveNumber() int {
	a := strings.Split(p.FEN, " ")
	m, _ := strconv.Atoi(a[5])
	return m
}

func (p *Position) ActiveColor() piece.Color {
	a := strings.Split(p.FEN, " ")
	var c piece.Color
	switch strings.ToLower(a[1]) {
	default:
		c = piece.White
	case "b":
		c = piece.Black
	}
	return c
}
