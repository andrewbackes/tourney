package models

import (
	"github.com/andrewbackes/chess/game"
	"github.com/andrewbackes/chess/piece"
)

type EndingCondition string

const (
	Threefold            EndingCondition = "Threefold Repetition"
	FiftyMoveRule                        = "Fifty-move Rule"
	Stalemate                            = "Stalemate"
	InsufficientMaterial                 = "Insufficient Material"
	Checkmate                            = "Checkmate"
	OutOfTime                            = "Out of Time"
	Resignation                          = "Resignation"
	IllegalMove                          = "Illegal Move"
	Error                                = "Error"
)

type Game struct {
	Id              Id                     `json:"id"`
	TournamentId    Id                     `json:"tournamentId"`
	Round           int                    `json:"round"`
	Status          Status                 `json:"status"`
	Result          Result                 `json:"result"`
	EndingCondition EndingCondition        `json:"endingCondition"`
	Contestants     map[piece.Color]Engine `json:"contestants"`
	TimeControl     game.TimeControl       `json:"timeControl"`
	Positions       []Position             `json:"positions"`
}

type CollapsedGame struct {
	Id              Id                     `json:"id"`
	TournamentId    Id                     `json:"tournamentId"`
	Round           int                    `json:"round"`
	Status          Status                 `json:"status"`
	Result          Result                 `json:"result"`
	EndingCondition EndingCondition        `json:"endingCondition"`
	Contestants     map[piece.Color]Engine `json:"contestants"`
	TimeControl     game.TimeControl       `json:"timeControl"`
}

func CollapseGame(g *Game) *CollapsedGame {
	return &CollapsedGame{
		Id:              g.Id,
		TournamentId:    g.TournamentId,
		Round:           g.Round,
		Status:          g.Status,
		Result:          g.Result,
		EndingCondition: g.EndingCondition,
		Contestants:     g.Contestants,
		TimeControl:     g.TimeControl,
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
