package model

import (
	"errors"
	"github.com/andrewbackes/tourney/model/structures"
	"gopkg.in/mgo.v2/bson"
	"sort"
	"time"
)

type tournamentList []*structures.Tournament

func (t tournamentList) Len() int           { return len(t) }
func (t tournamentList) Swap(i, j int)      { t[i], t[j] = t[j], t[i] }
func (t tournamentList) Less(i, j int) bool { return t[i].Created.Sub(t[j].Created) < 0 }

func (m *Model) AddTournament(t *structures.Tournament) (bson.ObjectId, error) {
	t.Init()
	t.Created = time.Now()
	m.tournamentMutex.Lock()
	m.tournaments[t.Id] = t
	m.tournamentMutex.Unlock()
	return t.Id, nil
}

func (m *Model) GetTournaments(filters map[string][]string) []*structures.Tournament {
	arr := make([]*structures.Tournament, 0, 0)
	filterIncompleteOnly := false
	if filters != nil {
		if v, exists := filters["completed"]; exists && len(v) > 0 && v[0] == "false" {
			filterIncompleteOnly = true
		}
	}
	m.tournamentMutex.RLock()
	for _, v := range m.tournaments {
		if !filterIncompleteOnly || (filterIncompleteOnly && !v.Complete()) {
			arr = append(arr, v)
		}
	}
	m.tournamentMutex.RUnlock()
	sort.Sort(tournamentList(arr))
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
