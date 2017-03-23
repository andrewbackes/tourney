package structures

import (
	"github.com/andrewbackes/chess/game"
	"gopkg.in/mgo.v2/bson"
	"strconv"
)

type Tournament struct {
	Id          bson.ObjectId                `json:"id" bson:"_id"`
	Tags        map[string]string            `json:"tags,omitempty" bson:"tags,omitempty"`
	TestSeats   int                          `json:"testSeats" bson:"testSeats"`
	Carousel    bool                         `json:"carousel" bson:"carousel"`
	Rounds      int                          `json:"rounds" bson:"rounds"`
	Contestants []Engine                     `json:"contestants" bson:"contestants"`
	Games       map[bson.ObjectId]*game.Game `json:"-" bson:"-"`
	queue       chan *game.Game
}

func NewTournament() *Tournament {
	return &Tournament{
		Id:          bson.NewObjectId(),
		Tags:        make(map[string]string),
		TestSeats:   1,
		Carousel:    true,
		Rounds:      0,
		Contestants: make([]Engine, 0),
		Games:       make(map[bson.ObjectId]*game.Game),
		queue:       make(chan *game.Game),
	}
}
func NewGameQueue(t *Tournament) chan *game.Game {
	queue := make(chan *game.Game, len(t.Games))
	ordered := make([]*game.Game, len(t.Games))
	for k := range t.Games {
		if r, err := strconv.Atoi(t.Games[k].Tags["round"]); err == nil {
			ordered[r] = t.Games[k]
		} else {
			panic(err)
		}
	}
	for _, v := range ordered {
		if result, exists := v.Tags["result"]; !exists || result == "*" {
			queue <- v
		}
	}
	return queue
}

func (t *Tournament) NextGame() *game.Game {
	select {
	case g := <-t.queue:
		return g
	default:
		return nil
	}
}

func (t *Tournament) Init() {
	if t.Id == "" {
		t.Id = bson.NewObjectId()
	}
	t.Contestants = IdentifyContestants(t)
	t.Games = NewGameList(t)
	t.queue = NewGameQueue(t)
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

func NewGameList(t *Tournament) map[bson.ObjectId]*game.Game {
	games := make(map[bson.ObjectId]*game.Game)

	round := 0
	for i := 0; i < t.TestSeats; i++ {
		if t.Carousel {
			for r := 0; r < t.Rounds; r = r + []int{2, 1}[t.Rounds%2] {
				for e := i + 1; e < len(t.Contestants); e++ {
					round++
					g := game.New()
					gid := bson.NewObjectId()
					g.Tags["round"] = strconv.Itoa(round)
					g.Tags["tournamentId"] = t.Id.Hex()
					g.Tags["id"] = gid.Hex()

					if r%2 == 0 {
						g.Tags["WhiteId"] = t.Contestants[i].Id.Hex()
						g.Tags["BlackId"] = t.Contestants[e].Id.Hex()
					} else {
						g.Tags["BlackId"] = t.Contestants[i].Id.Hex()
						g.Tags["WhiteId"] = t.Contestants[e].Id.Hex()
					}
					games[gid] = g

					if t.Rounds%2 == 0 {
						round++
						ng := game.New()
						ngid := bson.NewObjectId()
						ng.Tags["round"] = strconv.Itoa(round)
						ng.Tags["tournamentId"] = t.Id.Hex()
						ng.Tags["id"] = ngid.Hex()
						if r%2 == 0 {
							ng.Tags["WhiteId"] = t.Contestants[e].Id.Hex()
							ng.Tags["BlackId"] = t.Contestants[i].Id.Hex()
						} else {
							ng.Tags["BlackId"] = t.Contestants[i].Id.Hex()
							ng.Tags["WhiteId"] = t.Contestants[e].Id.Hex()
						}
						games[ngid] = ng
					}
				}
			}
		} else {
			// Non-Carousel:
			for e := i + 1; e < len(t.Contestants); e++ {
				//Now go around each opponent for that test seat:
				for r := 0; r < t.Rounds; r++ {
					round++
					g := game.New()
					gid := bson.NewObjectId()
					g.Tags["round"] = strconv.Itoa(round)
					g.Tags["tournamentId"] = t.Id.Hex()
					g.Tags["id"] = gid.Hex()
					if r%2 == 0 {
						g.Tags["WhiteId"] = t.Contestants[i].Id.Hex()
						g.Tags["BlackId"] = t.Contestants[e].Id.Hex()
					} else {
						g.Tags["BlackId"] = t.Contestants[i].Id.Hex()
						g.Tags["WhiteId"] = t.Contestants[e].Id.Hex()
					}
					games[gid] = g
				}
			}
		}
	}
	return games
}
