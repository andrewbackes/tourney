package structures

import (
	"gopkg.in/mgo.v2/bson"
	"sync"
)

// Game stores a game's state throughout its lifetime. It does not provide
// methods for playing a game. Playing a game is done through chess.Game.
type Game struct {
	sync.RWMutex
	Id        bson.ObjectId     `json:"id,omitempty" bson:"_id"`
	Tags      map[string]string `json:"tags,omitempty" bson:"tags"`
	Positions []*Position       `json:"positions,omitempty" bson:"positions"`
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

func (g *Game) Complete() bool {
	if result, exists := g.Tags["result"]; exists {
		return result == "1/2-1/2" || result == "1-0" || result == "0-1"
	}
	return false
}

func (g *Game) UpdateTags(t map[string]string) {
	if t != nil {
		for k, v := range t {
			g.Lock()
			g.Tags[k] = v
			g.Unlock()
		}
	}
}

func (g *Game) TournamentId() bson.ObjectId {
	g.RLock()
	tid := g.Tags["tournamentId"]
	g.RUnlock()
	return bson.ObjectIdHex(tid)
}
