package model

import (
	"errors"
	"github.com/andrewbackes/tourney/model/structures"
	"gopkg.in/mgo.v2/bson"
)

func (m *Model) AddTournament(t *structures.Tournament) (bson.ObjectId, error) {
	t.Init()
	m.tournamentMutex.Lock()
	m.tournaments[t.Id] = t
	m.tournamentMutex.Unlock()
	return t.Id, nil
}

func (m *Model) GetTournaments() []*structures.Tournament {
	arr := make([]*structures.Tournament, 0, len(m.tournaments))
	m.tournamentMutex.RLock()
	for _, v := range m.tournaments {
		arr = append(arr, v)
	}
	m.tournamentMutex.RUnlock()
	return arr
}

func (m *Model) GetTournament(id bson.ObjectId) (*structures.Tournament, error) {
	m.tournamentMutex.RLock()
	t, exists := m.tournaments[id]
	if !exists {
		return t, errors.New("Tournament not found")
	}
	m.tournamentMutex.RUnlock()
	return t, nil
}

func (m *Model) DeleteTournament(id bson.ObjectId) {
	m.tournamentMutex.Lock()
	delete(m.tournaments, id)
	m.tournamentMutex.Unlock()
}
