package main

import (
	"fmt"
	"github.com/andrewbackes/tourney/api"
	"github.com/andrewbackes/tourney/model"
	"github.com/andrewbackes/tourney/model/data/mongodb"
)

func main() {
	fmt.Println("Tourney")
	dao := mongodb.New("localhost")
	defer dao.Close()
	model := model.New(dao)
	server := api.New(model)
	server.ListenAndServe()
}
