package model

import (
	"gopkg.in/mgo.v2/bson"
)

type DataAccessor interface {
	GetTournaments() []Tournament
	AddTournament(Tournament)
	DeleteTournament(id bson.ObjectId)
	Close()
}
