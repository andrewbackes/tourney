/*******************************************************************************

 Project: Tourney

 Module: worker
 Created: 12/3/2014
 Author(s): Andrew Backes
 Description:

 TODO:
 	-

*******************************************************************************/

package main

import (
	"fmt"
	"net"
	"net/rpc"
)

// Wrapper for functions to be played with net/rpc :
type Worker struct{}

func (W *Worker) PlayGame(G Game, CompletedGame *Game) error {
	fmt.Println("Server says to play a game.")

	fmt.Println(G)
	G.PrintHUD()

	WorkingGame := G

	PlayGame(&WorkingGame)

	*CompletedGame = WorkingGame
	fmt.Println("Done Playing game.")

	fmt.Println(WorkingGame)
	WorkingGame.PrintHUD()

	return nil
}

func ConnectAndWait(address string) {
	// First connect to the host:
	fmt.Print("\nConnecting to " + address + "... ")
	conn, err := net.Dial("tcp", address)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	fmt.Println("Success.")

	fmt.Println("Waiting on server...")
	rpc.Register(new(Worker))
	rpc.ServeConn(conn) //don't forget that this blocks

	fmt.Println("Connection closed.")
}
