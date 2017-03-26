package main

import (
	"fmt"
	"github.com/andrewbackes/tourney/api"
	"github.com/andrewbackes/tourney/model"
)

func main() {
	fmt.Println("Tourney")
	//dao := mongodb.New("localhost")
	//defer dao.Close()
	model := model.New()
	server := api.New(model)
	server.ListenAndServe()
}
