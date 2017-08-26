package models

import (
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
	"github.com/andrewbackes/chess/position/move"
	"time"
)

type Game struct {
	Id           Id                     `json:"id"`
	TournamentId Id                     `json:"tournamentId"`
	Status       Status                 `json:"status"`
	Contestants  map[piece.Color]Engine `json:"contestants"`
	TimeControl  game.TimeControl       `json:"timeControl"`
	Positions    []Position             `json:"positions"`
}

type Position struct {
	FEN       string                        `json:"fen"`
	LastMove  move.Move                     `json:"lastMove"`
	MovesLeft map[piece.Color]int           `json:"movesLeft"`
	Clocks    map[piece.Color]time.Duration `json:"clock"`
}

func NewGame(tid Id, c game.TimeControl, w, b Engine) *Game {
	return &Game{
		Id:           NewId(),
		TournamentId: tid,
		TimeControl:  c,
		Contestants: map[piece.Color]Engine{
			piece.White: w,
			piece.Black: b,
		},
		Positions: make([]Position, 0),
	}
}
