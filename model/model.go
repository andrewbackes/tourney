package model

import (
	"errors"
	"github.com/andrewbackes/tourney/model/data"
	"github.com/andrewbackes/tourney/model/structures"
	"gopkg.in/mgo.v2/bson"
)

var (
	// ErrorNotFound Resource can not be found.
	ErrorNotFound = errors.New("Not found")
)

const (
	TournamentQueueBuffer = 1000
)

type Model struct {
	Tournaments map[bson.ObjectId]*structures.Tournament
	Engines     map[bson.ObjectId]*structures.Engine
	Books       map[bson.ObjectId]*structures.Book
	Workers     map[bson.ObjectId]*structures.Worker

	queue chan *structures.Tournament
	dao   data.Accessor
}

func New(dao data.Accessor) *Model {
	m := Model{
		Tournaments: make(map[bson.ObjectId]*structures.Tournament),
		queue:       make(chan *structures.Tournament, TournamentQueueBuffer),
		dao:         dao,
	}
	/*
		ts := m.dao.GetTournaments()
		for _, v := range ts {
			m.Tournaments[v.Id] = &v
		}
	*/
	return &m
}
