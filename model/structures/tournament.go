package structures

import (
	"gopkg.in/mgo.v2/bson"
)

type Tournament struct {
	Id          bson.ObjectId     `json:"id" bson:"_id"`
	Tags        map[string]string `json:"tags,omitempty" bson:"tags,omitempty"`
	TestSeats   int               `json:"testSeats" bson:"testSeats"`
	Carousel    bool              `json:"carousel" bson:"carousel"`
	Rounds      int               `json:"rounds" bson:"rounds"`
	Contestants []Engine          `json:"contestants" bson:"contestants"`
}

func (t *Tournament) IdentifyContestants() []Engine {
	engines := make([]Engine, 0, len(t.Contestants))
	for _, e := range t.Contestants {
		if e.Id == bson.ObjectId("") {
			e.Id = bson.NewObjectId()
		}
		engines = append(engines, e)
	}
	return engines
}

func (t *Tournament) GenerateGames() []Game {
	var games []Game
	for i := 0; i < t.TestSeats; i++ {
		if t.Carousel {
			for r := 0; r < t.Rounds; r = r + []int{2, 1}[t.Rounds%2] {
				for e := i + 1; e < len(t.Contestants); e++ {
					g := Game{
						Id:           bson.NewObjectId(),
						TournamentId: t.Id,
					}
					if r%2 == 0 {
						g.WhiteId = t.Contestants[i].Id
						g.BlackId = t.Contestants[e].Id
					} else {
						g.BlackId = t.Contestants[i].Id
						g.WhiteId = t.Contestants[e].Id
					}
					games = append(games, g)
					if t.Rounds%2 == 0 {
						ng := Game{
							Id:           bson.NewObjectId(),
							TournamentId: t.Id,
						}
						if r%2 == 0 {
							ng.WhiteId = t.Contestants[e].Id
							ng.BlackId = t.Contestants[i].Id
						} else {
							ng.BlackId = t.Contestants[i].Id
							ng.WhiteId = t.Contestants[e].Id
						}
						games = append(games, ng)
					}
				}
			}
		} else {
			// Non-Carousel:
			for e := i + 1; e < len(t.Contestants); e++ {
				//Now go around each opponent for that test seat:
				for r := 0; r < t.Rounds; r++ {
					g := Game{
						Id:           bson.NewObjectId(),
						TournamentId: t.Id,
					}
					if r%2 == 0 {
						g.WhiteId = t.Contestants[i].Id
						g.BlackId = t.Contestants[e].Id
					} else {
						g.BlackId = t.Contestants[i].Id
						g.WhiteId = t.Contestants[e].Id
					}
					games = append(games, g)
				}
			}
		}
	}
	return games
}
