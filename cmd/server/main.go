package main

import (
	"fmt"
	"github.com/andrewbackes/tourney/db/memdb"
	"github.com/andrewbackes/tourney/server"
)

func main() {
	fmt.Println("Server")
	db := memdb.NewMemDB()
	s := server.NewServer(":9090", db)
	s.Start()
}
