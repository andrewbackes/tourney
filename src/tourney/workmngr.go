/*******************************************************************************

 Project: Tourney

 Module: workmanager
 Created: 12/3/2014
 Author(s): Andrew Backes
 Description:

 TODO:
 	-customizable ports

*******************************************************************************/

package main

import (
	"fmt"
	"net"
	"net/http"
	"net/rpc"
)

type WorkManager struct {
	ConnectedWorkers map[*Worker]struct{}
	WorkerQue        chan *Worker
	CompletedGames   chan Game
	//TODO: keep track of what is assigned to which worker
}

func NewWorkManager(T *Tourney) *WorkManager {
	M := &WorkManager{
		ConnectedWorkers: make(map[*Worker]struct{}),
		WorkerQue:        make(chan *Worker, len(T.GameList)), // BUG: possible deadlocks! think about what happens if more workers connect than games! this will cause a deadlock somewhere!
		CompletedGames:   make(chan Game, len(T.GameList)),
	}
	return M
}

func (M *WorkManager) ConnectWorker(conn net.Conn) {
	fmt.Println(conn.RemoteAddr(), "connected.")
	//rpcConn := rpc.NewClient(conn) // don't forget this is a pointer
	W := &Worker{
		Address: conn.RemoteAddr(),
		RPC:     rpc.NewClient(conn),
	}
	M.ConnectedWorkers[W] = struct{}{}
	M.WorkerQue <- W
	fmt.Println("Worker added to que.")
}

func (M *WorkManager) DisconnectWorker(W *Worker) {
	//Close the connection:
	W.RPC.Close()
	// Remove from the connected workers list:
	delete(M.ConnectedWorkers, W)
}

func (M *WorkManager) DisconnectAll() {
	fmt.Println("Disconnecting all workers.")
	for key, _ := range M.ConnectedWorkers {
		M.DisconnectWorker(key)
	}
}

func (M *WorkManager) ListenForWorkers() {
	// Setup Server:
	fmt.Println("Waiting for workers on port 9000...")
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

func (M *WorkManager) RemotelyPlayGame(W *Worker, GameToPlay Game) {
	//TODO: are there 3 copies of the game at this point!? fix this!

	fmt.Println("Round", GameToPlay.Round, "being played by", W.Address)
	var CompletedGame Game
	GameToPlay.Site = fmt.Sprint(W.Address)

	// Make sure MD5 sums are set:
	GameToPlay.Player[0].ValidateEngineFile()
	GameToPlay.Player[1].ValidateEngineFile()

	// Start playing:
	err := W.RPC.Call("Worker.PlayGame", GameToPlay, &CompletedGame)
	if err != nil {
		fmt.Println("Error remotely playing game:", err)
		return
	}
	fmt.Println("Round", CompletedGame.Round, "completed.")

	// Que the completed game for reconsiliation:
	M.CompletedGames <- CompletedGame
	fmt.Println(W.Address, "added to the worker que.")
	M.WorkerQue <- W
}

func SyncGames(T *Tourney, M *WorkManager) {
	// Note: 	This should NEVER be put in its own goroutine.
	//			It would most definitely cause a race condition somewhere.
GAMESYNC:
	for {
		select {
		case GameToUpdate := <-M.CompletedGames:
			fmt.Println("Syncronizing round", GameToUpdate.Round)
			T.GameList[GameToUpdate.Round-1] = GameToUpdate
			// Save progress:
			if err := Save(T); err != nil {
				fmt.Println(err)
			}
		default:
			//fmt.Println("No games to sync.")
			break GAMESYNC
		}
	}
}

func ServeEngineFiles() {
	// TODO: customizable ports
	FilePath := "/Users/Andrew/Documents/"
	fmt.Println("Serving game engines on port 9001")
	http.ListenAndServe(":9001", http.FileServer(http.Dir(FilePath)))
}

func HostTourney(T *Tourney) error {

	fmt.Println("\n\nHosting:", T.Event)
	M := NewWorkManager(T)
	go ServeEngineFiles() // TODO: BUG: race condition here. if server isnt up and clients are trying to download the files.
	go M.ListenForWorkers()

	// TODO: refactoring required: consolidate with RunTourney()

	// Helper function:
	var GameIndex, IterationCount int
	GameIndex = -1                // hack so that the first returned value is 0
	var NextGameIndex func() *int // prototype
	NextGameIndex = func() *int {
		// gets the index of the next game that needs to be played.
		// goes through the gamelist in order assigning indexes.
		// then loops back around in case some games didnt get returned.
		GameIndex = (GameIndex + 1) % len(T.GameList)
		if T.GameList[GameIndex].Completed {
			IterationCount++
			if IterationCount > len(T.GameList) {
				return nil
			}
			return NextGameIndex()
		}
		IterationCount = 0
		return &GameIndex
	}

	// Wait for a free worker and assign that worker the next game. Do this until all work is complete:
	WorkComplete := false
	for !WorkComplete {
		select {
		case <-T.Done:
			// user force quits
			WorkComplete = true //hack
			break
		case freeWorker := <-M.WorkerQue:
			// Since a game is que'd as complete before a worker is que'd as free, we need to sync a game:
			SyncGames(T, M)
			// Figure out what game to assign this free worker. If any.
			if pNextGameIndex := NextGameIndex(); pNextGameIndex != nil {
				fmt.Println("Round", *pNextGameIndex+1, "started.")
				// play the opening.
				fmt.Print("Playing from opening book... ")
				// NOTE: if the opening was already played, this should not error, but just continue on like normal.
				if err := PlayOpening(T, *pNextGameIndex); err != nil {
					fmt.Println("Failed:", err.Error())
					M.WorkerQue <- freeWorker // BUG: This can deadlock!
					// 		what if there are N games to play and this is the N+1 worker.
					// 		The only reciever of this channel is in the select that this is nested in.
					break
				}
				fmt.Println("Success.")
				// Remotely play game:
				GameToPlay := T.GameList[*pNextGameIndex] // make a copy to prevent race conditions.
				go M.RemotelyPlayGame(freeWorker, GameToPlay)

			} else {
				//done with the tourey!
				fmt.Println("Done hosting Tourney.")
				//M.WorkerQue <- freeWorker
				WorkComplete = true
			}
		}
	}

	// Show results:
	fmt.Print(SummarizeResults(T))

	M.DisconnectAll()
	return nil
}
