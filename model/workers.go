package model

import (
	"github.com/andrewbackes/tourney/model/structures"
	"gopkg.in/mgo.v2/bson"
)

// AddWorker returns a Game Id for that worker to play.
func (m *Model) AddWorker(w *structures.Worker) bson.ObjectId {
	m.Workers[w.Id] = w
	// TODO: request new game
	return w.GameId
}
