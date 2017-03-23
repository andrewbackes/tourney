package model

import (
	"errors"
	"github.com/andrewbackes/tourney/model/structures"
	"gopkg.in/mgo.v2/bson"
)

func (m *Model) AddTournament(t *structures.Tournament) (bson.ObjectId, error) {
	t.Init()
	m.queue <- t
	return t.Id, nil
}

func (m *Model) GetTournaments() []structures.Tournament {
	arr := make([]structures.Tournament, 0, len(m.Tournaments))
	for _, v := range m.Tournaments {
		arr = append(arr, *v)
	}
	return arr
}

func (m *Model) GetTournament(id bson.ObjectId) (structures.Tournament, error) {
	t, exists := m.Tournaments[id]
	if !exists {
		return *t, errors.New("Tournament not found")
	}
	return *t, nil
}

func (m *Model) DeleteTournament(id bson.ObjectId) {
	delete(m.Tournaments, id)
}
