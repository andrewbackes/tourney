package main

import (
	"flag"
	"fmt"
	"github.com/andrewbackes/tourney/api"
	"github.com/andrewbackes/tourney/client"
	"github.com/andrewbackes/tourney/model"
)

func main() {
	fmt.Println("Tourney")
	connect := flag.String("connect", "", "Connect to a tourney server.")
	flag.Parse()
	if *connect != "" {
		fmt.Println("Connecting to ", *connect)
		c := client.New(*connect)
		c.ConnectAndPlay()
	} else {
		fmt.Println("Running on port 8080")
		//dao := mongodb.New("localhost")
		//defer dao.Close()
		model := model.New()
		server := api.New(model)
		server.ListenAndServe()
	}
}
