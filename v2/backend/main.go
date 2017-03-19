package main

import (
	"encoding/json"
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

/*
type InMemoryStorage struct{}
func (m InMemoryStorage) Close()                            {}
func (m InMemoryStorage) GetTournaments() []Tournament      { return make([]Tournament, 0) }
func (m InMemoryStorage) AddTournament(t Tournament)        {}
func (m InMemoryStorage) DeleteTournament(id bson.ObjectId) {}
*/

func main() {
	fmt.Println("Tourney")
	model := NewModel("mongodb", "localhost")
	defer model.dao.Close()
	http.HandleFunc("/api/v2/tournaments", func(w http.ResponseWriter, req *http.Request) {
		switch req.Method {
		case http.MethodPost:
			decoder := json.NewDecoder(req.Body)
			var t Tournament
			err := decoder.Decode(&t)
			if err != nil {
				panic(err)
			}
			defer req.Body.Close()
			id := model.CreateTournament(t)
			w.Write([]byte(fmt.Sprintf("{\"id\":\"%s\"}", id.Hex())))
		case http.MethodGet:
			t := model.GetTournaments()
			encoder := json.NewEncoder(w)
			encoder.SetIndent("", "  ")
			encoder.Encode(t)
		}
	})
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Println(err)
	}
}
