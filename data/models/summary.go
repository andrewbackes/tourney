package models

import (
	"github.com/andrewbackes/chess/piece"
)

type Summary struct {
	Complete   int          `json:"complete"`
	Incomplete int          `json:"incomplete"`
	Stats      map[Id]Stats `json:"stats"`
}

type Stats struct {
	Wins       int `json:"wins"`
	Losses     int `json:"losses"`
	Draws      int `json:"draws"`
	Incomplete int `json:"incomplete"`
}

func NewSummary(contestants []Engine, gl []*Game) *Summary {
	s := Summary{
		Complete:   0,
		Incomplete: 0,
		Stats:      make(map[Id]Stats),
	}
	for _, g := range gl {
		w := s.Stats[g.Contestants[piece.White].Id]
		b := s.Stats[g.Contestants[piece.Black].Id]
		switch g.Result {
		case Incomplete:
			s.Incomplete++
			w.Incomplete++
			b.Incomplete++
		case White:
			s.Complete++
			w.Wins++
			b.Losses++
		case Black:
			s.Complete++
			b.Wins++
			w.Losses++
		case Draw:
			s.Complete++
			w.Draws++
			b.Draws++
		}
		s.Stats[g.Contestants[piece.White].Id] = w
		s.Stats[g.Contestants[piece.Black].Id] = b
	}
	return &s
}
