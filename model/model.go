package model

import (
	"errors"
	"github.com/andrewbackes/tourney/model/data"
	"github.com/andrewbackes/tourney/model/structures"
	"gopkg.in/mgo.v2/bson"
	"sync"
)

var (
	// ErrorNotFound Resource can not be found.
	ErrorNotFound = errors.New("Not found")
)

const (
	TournamentQueueBuffer = 1000
)

type Model struct {
	tournaments     map[bson.ObjectId]*structures.Tournament
	tournamentMutex sync.RWMutex

	Engines map[bson.ObjectId]*structures.Engine
	Books   map[bson.ObjectId]*structures.Book
	Workers map[bson.ObjectId]*structures.Worker

	dao  data.Accessor
	done chan struct{}
}

func New() *Model {
	m := Model{
		tournaments: make(map[bson.ObjectId]*structures.Tournament),
		done:        make(chan struct{}),
	}
	/*
		ts := m.dao.GetTournaments()
		for _, v := range ts {
			m.Tournaments[v.Id] = &v
		}
	*/

	return &m
}

func (m *Model) Stop() {
	close(m.done)
}
