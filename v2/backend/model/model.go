package model

import (
	"gopkg.in/mgo.v2/bson"
)

type Model struct {
	Tournaments map[bson.ObjectId]Tournament
	dao         DataAccessor
}

func NewModel(dao string, args ...string) *Model {
	m := Model{
		Tournaments: make(map[bson.ObjectId]Tournament),
	}
	switch dao {
	case "mongodb":
		url := "localhost"
		if len(args) > 0 {
			url = args[0]
		}
		m.dao = NewMongoDB(url)
	case "inmemory":
		//m.dao = InMemoryStorage{}
	default:
		panic("Invalid persister " + dao)
	}
	ts := m.dao.GetTournaments()
	for _, v := range ts {
		m.Tournaments[v.Id] = v
	}
	return &m
}

func (m *Model) CreateTournament(t Tournament) bson.ObjectId {
	t.Id = bson.NewObjectId()
	m.Tournaments[t.Id] = t
	m.dao.AddTournament(t)
	return t.Id
}

func (m *Model) GetTournaments() []Tournament {
	arr := make([]Tournament, 0, len(m.Tournaments))
	for _, v := range m.Tournaments {
		arr = append(arr, v)
	}
	return arr
}

func (m *Model) DeleteTournament(id bson.ObjectId) {
	delete(m.Tournaments, id)
	m.dao.DeleteTournament(id)
}
