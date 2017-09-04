package models

import (
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
)

type Game struct {
	Id           Id                     `json:"id"`
	TournamentId Id                     `json:"tournamentId"`
	Status       Status                 `json:"status"`
	Contestants  map[piece.Color]Engine `json:"contestants"`
	TimeControl  game.TimeControl       `json:"timeControl"`
	Positions    []Position             `json:"positions"`
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
		Positions: []Position{StartPosition(c)},
	}
}
