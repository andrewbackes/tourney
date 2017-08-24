package structures

import (
	"fmt"
	"gopkg.in/mgo.v2/bson"
	"strconv"
	"sync"
	"time"
)

type Tournament struct {
	sync.RWMutex
	Id          bson.ObjectId           `json:"id" bson:"_id"`
	Created     time.Time               `json:"created,omitempty" bson:"created,omitempty"`
	Tags        map[string]string       `json:"tags,omitempty" bson:"tags,omitempty"`
	TestSeats   int                     `json:"testSeats" bson:"testSeats"`
	Carousel    bool                    `json:"carousel" bson:"carousel"`
	Rounds      int                     `json:"rounds" bson:"rounds"`
	Contestants []Engine                `json:"contestants" bson:"contestants"`
	games       map[bson.ObjectId]*Game `json:"-" bson:"-"`
	queue       chan *Game
}

func NewTournament() *Tournament {
	return &Tournament{
		Id:          bson.NewObjectId(),
		Tags:        make(map[string]string),
		TestSeats:   1,
		Carousel:    true,
		Rounds:      0,
		Contestants: make([]Engine, 0),
		games:       make(map[bson.ObjectId]*Game),
		queue:       make(chan *Game),
	}
}

func NewGameQueue(t *Tournament) chan *Game {
	queue := make(chan *Game, len(t.games))
	ordered := make([]*Game, len(t.games))
	for k := range t.games {
		if r, err := strconv.Atoi(t.games[k].Tags["round"]); err == nil {
			ordered[r-1] = t.games[k]
		} else {
			panic(err)
		}
	}
	for _, v := range ordered {
		if !v.Complete() {
			queue <- v
		}
	}
	return queue
}

func (t *Tournament) NextGame() *Game {
	if len(t.queue) == 0 && !t.Complete() {
		t.Lock()
		t.queue = NewGameQueue(t)
		t.Unlock()
	}
	select {
	case g := <-t.queue:
		return g
	default:
		fmt.Printf("Game queue depleted for tournament %s\n", t.Id.Hex())
		return nil
	}
}

func (t *Tournament) Init() {
	t.Lock()
	if t.Id == "" {
		t.Id = bson.NewObjectId()
	}
	t.Contestants = IdentifyContestants(t)
	t.games = NewGameList(t)
	t.queue = NewGameQueue(t)
	t.Unlock()
}

func IdentifyContestants(t *Tournament) []Engine {
	engines := make([]Engine, 0, len(t.Contestants))
	for _, e := range t.Contestants {
		if e.Id == bson.ObjectId("") {
			e.Id = bson.NewObjectId()
		}
		engines = append(engines, e)
	}
	return engines
}

func NewGameList(t *Tournament) map[bson.ObjectId]*Game {
	games := make(map[bson.ObjectId]*Game)
	round := 0
	for i := 0; i < t.TestSeats; i++ {
		if t.Carousel {
			for r := 0; r < t.Rounds; r = r + []int{2, 1}[t.Rounds%2] {
				for e := i + 1; e < len(t.Contestants); e++ {
					round++
					g := NewGame()
					g.Tags["round"] = strconv.Itoa(round)
					g.Tags["tournamentId"] = t.Id.Hex()
					if r%2 == 0 {
						g.Tags["whiteId"] = t.Contestants[i].Id.Hex()
						g.Tags["blackId"] = t.Contestants[e].Id.Hex()
					} else {
						g.Tags["blackId"] = t.Contestants[i].Id.Hex()
						g.Tags["whiteId"] = t.Contestants[e].Id.Hex()
					}
					games[g.Id] = g
					if t.Rounds%2 == 0 {
						round++
						ng := NewGame()
						ng.Tags["round"] = strconv.Itoa(round)
						ng.Tags["tournamentId"] = t.Id.Hex()

						if r%2 == 0 {
							ng.Tags["whiteId"] = t.Contestants[e].Id.Hex()
							ng.Tags["blackId"] = t.Contestants[i].Id.Hex()
						} else {
							ng.Tags["blackId"] = t.Contestants[i].Id.Hex()
							ng.Tags["whiteId"] = t.Contestants[e].Id.Hex()
						}
						games[ng.Id] = ng
					}
				}
			}
		} else {
			// Non-Carousel:
			for e := i + 1; e < len(t.Contestants); e++ {
				//Now go around each opponent for that test seat:
				for r := 0; r < t.Rounds; r++ {
					round++
					g := NewGame()
					g.Tags["round"] = strconv.Itoa(round)
					g.Tags["tournamentId"] = t.Id.Hex()
					if r%2 == 0 {
						g.Tags["whiteId"] = t.Contestants[i].Id.Hex()
						g.Tags["blackId"] = t.Contestants[e].Id.Hex()
					} else {
						g.Tags["blackId"] = t.Contestants[i].Id.Hex()
						g.Tags["whiteId"] = t.Contestants[e].Id.Hex()
					}
					games[g.Id] = g
				}
			}
		}
	}
	return games
}

func (t *Tournament) GetGame(id bson.ObjectId) *Game {
	t.RLock()
	g, exists := t.games[id]
	t.RUnlock()
	if exists {
		return g
	}
	panic("gid not found")
}

func (t *Tournament) GetGames() []*Game {
	arr := make([]*Game, 0, len(t.games))
	t.RLock()
	for _, v := range t.games {
		arr = append(arr, v)
	}
	t.RUnlock()
	return arr
}

func (t *Tournament) Complete() bool {
	t.RLock()
	for k, _ := range t.games {
		if !t.games[k].Complete() {
			return false
		}
	}
	t.RUnlock()
	return true
}
