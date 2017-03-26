package structures

import (
	"gopkg.in/mgo.v2/bson"
)

type Worker struct {
	Id           bson.ObjectId
	GameId       bson.ObjectId
	TournamentId bson.ObjectId
}

func NewWorker() *Worker {
	return &Worker{
		Id: bson.NewObjectId(),
	}
}
