package structures

import (
	"github.com/andrewbackes/chess/fen"
	"github.com/andrewbackes/chess/position"
)

// Position represents the game's state after the associated move has been made.
type Position struct {
	position.Position
	Move position.Move `json:"move" bson:"move"`
	FEN  string        `json:"fen" bson:"fen"`
}

func NewPosition() *Position {
	p := position.New()
	p.MoveNumber = 0
	f, _ := fen.Encode(p)
	return &Position{
		*p,
		position.NullMove,
		f,
	}
}
