package models

import (
	"github.com/andrewbackes/chess/game"
)

// Settings are for configurating a tournament.
type Settings struct {
	TestSeats   int              `json:"testSeats"`
	Carousel    bool             `json:"carousel"`
	Rounds      int              `json:"rounds"`
	Engines     []Engine         `json:"engines"`
	TimeControl game.TimeControl `json:"timeControl"`
	Opening     Opening          `json:"opening"`
}

// Opening dictates how an opening book is used in a tournament.
type Opening struct {
	Book      Book `json:"book"`
	Depth     int  `json:"depth"`
	Randomize bool `json:"randomize"`
}
