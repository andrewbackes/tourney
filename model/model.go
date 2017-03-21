package model

import (
	"errors"
	"github.com/andrewbackes/tourney/model/data"
	"github.com/andrewbackes/tourney/model/structures"
	"gopkg.in/mgo.v2/bson"
)

type Model struct {
	Tournaments map[bson.ObjectId]structures.Tournament
	dao         data.Accessor
}

func New(dao data.Accessor) *Model {
	m := Model{
		Tournaments: make(map[bson.ObjectId]structures.Tournament),
		dao:         dao,
	}
	ts := m.dao.GetTournaments()
	for _, v := range ts {
		m.Tournaments[v.Id] = v
	}
	return &m
}

func (m *Model) CreateTournament(t structures.Tournament) bson.ObjectId {
	t.Id = bson.NewObjectId()
	m.Tournaments[t.Id] = t
	m.dao.AddTournament(t)
	return t.Id
}

func (m *Model) GetTournaments() []structures.Tournament {
	arr := make([]structures.Tournament, 0, len(m.Tournaments))
	for _, v := range m.Tournaments {
		arr = append(arr, v)
	}
	return arr
}

func (m *Model) GetTournament(id bson.ObjectId) (structures.Tournament, error) {
	t, exists := m.Tournaments[id]
	if !exists {
		return t, errors.New("Tournament not found")
	}
	return t, nil
}

func (m *Model) DeleteTournament(id bson.ObjectId) {
	delete(m.Tournaments, id)
	m.dao.DeleteTournament(id)
}
