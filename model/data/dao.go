package data

import (
	"github.com/andrewbackes/tourney/model/structures"
	"gopkg.in/mgo.v2/bson"
)

type Accessor interface {
	GetTournaments() []structures.Tournament
	AddTournament(*structures.Tournament)
	DeleteTournament(id bson.ObjectId)
	Close()
}
