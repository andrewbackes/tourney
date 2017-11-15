package models

import (
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
)

type Game struct {
	Id           Id                     `json:"id"`
	TournamentId Id                     `json:"tournamentId"`
	Status       Status                 `json:"status"`
	Result       Result                 `json:"result"`
	Contestants  map[piece.Color]Engine `json:"contestants"`
	TimeControl  game.TimeControl       `json:"timeControl"`
	Positions    []Position             `json:"positions"`
}

type CollapsedGame struct {
	Id           Id                     `json:"id"`
	TournamentId Id                     `json:"tournamentId"`
	Status       Status                 `json:"status"`
	Result       Result                 `json:"result"`
	Contestants  map[piece.Color]Engine `json:"contestants"`
	TimeControl  game.TimeControl       `json:"timeControl"`
}

func CollapseGame(g *Game) *CollapsedGame {
	return &CollapsedGame{
		Id:           g.Id,
		TournamentId: g.TournamentId,
		Status:       g.Status,
		Result:       g.Result,
		Contestants:  g.Contestants,
		TimeControl:  g.TimeControl,
	}
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
