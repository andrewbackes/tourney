/*******************************************************************************

 Project: Tourney

 Module: workmanager
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

type WorkManager struct {
	Workers   map[*rpc.Client]struct{} // Holds all connected workers.
	WorkerQue chan *rpc.Client
	//TODO: keep track of what is assigned to which worker
}

func NewWorkManager() *WorkManager {
	M := &WorkManager{
		Workers:   make(map[*rpc.Client]struct{}),
		WorkerQue: make(chan *rpc.Client),
	}
	return M
}

func (M *WorkManager) ConnectWorker(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "connected.")
	rpcConn := rpc.NewClient(conn) // don't forget this is a pointer
	M.Workers[rpcConn] = struct{}{}
	M.WorkerQue <- rpcConn
	fmt.Println("Worker added to que.")
}

func (M *WorkManager) ListenForWorkers() {
	// Setup Server:
	fmt.Println("\nListening on port 9000...")
	server, err := net.Listen("tcp", ":9000") //TODO: user chosen port.
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer server.Close()

	// Start listening:
	for {

		// Wait for a connection:
		conn, err := server.Accept() // TODO: add a timeout
		if err != nil {
			fmt.Println("Client connection error:", err.Error())
			continue
		}

		// Establish the incomming connection:
		go M.ConnectWorker(conn)
	}

}

func HostTourney(T *Tourney) error {

	fmt.Println("Hosting Tourney.")
	M := NewWorkManager()
	go M.ListenForWorkers()

	freeworker := <-M.WorkerQue

	game := T.GameList[20]
	fmt.Println(game)
	game.PrintHUD()

	var completedGame Game
	fmt.Println("Attempting to play game.")

	err := freeworker.Call("Worker.PlayGame", game, &completedGame)
	if err != nil {
		fmt.Println("Internal error:", err)
	}

	fmt.Println("Done playing game.")

	fmt.Println(completedGame)
	completedGame.PrintHUD()

	fmt.Println("Done hosting Tourney.")
	return nil
}
