package structures

import (
	"gopkg.in/mgo.v2/bson"
)

type Game struct {
	Id           bson.ObjectId     `json:"id" bson:"_id"`
	TournamentId bson.ObjectId     `json:"tournamentId" bson:"tournamentId"`
	WhiteId      bson.ObjectId     `json:"whiteId" bson:"whiteId"`
	BlackId      bson.ObjectId     `json:"blackId" bson:"blackId"`
	Tags         map[string]string `json:"tags,omitempty" bson:"tags,omitempty"`
}
