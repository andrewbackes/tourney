package structures

import (
	"gopkg.in/mgo.v2/bson"
	"sync"
)

// Game stores a game's state throughout its lifetime. It does not provide
// methods for playing a game. Playing a game is done through chess.Game.
type Game struct {
	sync.RWMutex
	Id        bson.ObjectId     `json:"id" bson:"_id"`
	Tags      map[string]string `json:"tags" bson:"tags"`
	Positions []*Position       `json:"positions" bson:"positions"`
}

func NewGame() *Game {
	return &Game{
		Id:        bson.NewObjectId(),
		Tags:      make(map[string]string),
		Positions: []*Position{NewPosition()},
	}
}

func (g *Game) GetPositions() []*Position {
	arr := make([]*Position, 0, len(g.Positions))
	g.RLock()
	for _, v := range g.Positions {
		arr = append(arr, v)
	}
	g.RUnlock()
	return arr
}

func (g *Game) GetPosition(index int) *Position {
	var p *Position
	g.RLock()
	p = g.Positions[index]
	g.RUnlock()
	return p
}

func (g *Game) AddPosition(p *Position) {
	index := 2*p.MoveNumber - (1 - int(p.ActiveColor))
	g.Lock()
	for len(g.Positions) < index+1 {
		g.Positions = append(g.Positions, nil)
	}
	g.Positions[index] = p
	g.Unlock()
}
