package mongodb

import (
	"github.com/andrewbackes/tourney/model/structures"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoDB struct {
	session *mgo.Session
	db      string
}

func New(url string) *MongoDB {
	session, err := mgo.Dial(url)
	if err != nil {
		panic(err)
	}
	return &MongoDB{
		db:      "tourney",
		session: session,
	}
}

func (m *MongoDB) Close() {
	m.session.Close()
}

func (m *MongoDB) DeleteTournament(id bson.ObjectId) {

}

func (m *MongoDB) GetTournaments() []structures.Tournament {
	s := m.session.Copy()
	c := s.DB(m.db).C("tournaments")
	var ts []structures.Tournament
	err := c.Find(nil).All(&ts)
	if err != nil {
		panic(err)
	}
	return ts
}

func (m *MongoDB) AddTournament(t structures.Tournament) {
	s := m.session.Copy()
	c := s.DB(m.db).C("tournaments")
	err := c.Insert(&t)
	if err != nil {
		panic(err)
	}
}
