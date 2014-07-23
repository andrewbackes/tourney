/*

 Project: Tourney

 Module: games
 Description: game object

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*/

package main

import (
	"fmt"
	"time"
	//"strings"
)

type Color uint8

const (
	WHITE   Color = 0
	BLACK   Color = 1
	BOTH    Color = 2
	NEITHER Color = 2
	DRAW    Color = 2
)

type Game struct {

	// Header info (General info, usually needed for pgn):
	time      int64
	moves     int64
	repeating bool
	player    [2]Engine //white=0,black=1

	//TODO: the rest of the PGN required heading. i think there are 8?

	// Control info:
	timer     [2]int64 //the game clock for each side. milliseconds.
	movesToGo int64    //moves until time control

	// Game run time info:
	state     Status //UNSTARTED,RUNNING,STOPPED
	toMove    Color  //WHITE or BLACK only
	board     Board
	moveList  []Move //move history
	completed bool   // TODO:  change this to reflect how the game ended. time, checkmate, adjunction, etc

	// Post game info:
	result Color //WHITE,BLACK,DRAW - should be set when game state is changed to STOPPED

}

/*******************************************************************************

	Primary function. Manages the game itself. Loops until instructed not to.

*******************************************************************************/

func (G *Game) playLoop() error {

	for G.state == RUNNING {

		for color := WHITE; color <= BLACK; color++ {
			//Tell the engine to set its board:
			if err := G.player[color].Set(G.moveList); err != nil { //TODO: should probably pass by ref
				G.Stop()
				break
			}

			//Request a move from the engine:
			startTime := time.Now()
			move, err := G.player[color].Move(G.timer, G.movesToGo)
			if err != nil {
				G.Stop()
				break
			}
			endTime := time.Now()
			lapsed := endTime.Sub(startTime)

			// Adjust time control:
			G.timer[color] -= int64(lapsed.Seconds() * 1000)
			if G.timer[color] <= 0 {
				G.GameOver(color)
			}

			G.MakeMove(move)

			fmt.Print([]string{"WHITE", "BLACK"}[color], "> ", move.algebraic, "\n")
		}

		G.movesToGo -= 1
		if G.movesToGo == 0 && G.repeating == true {
			G.resetTimeControl()
		}

		// temporarly only make one move each:
		//G.Stop()
	}
	return nil
}

/*******************************************************************************

	Functions that control the operating state of the game:

*******************************************************************************/

func (G *Game) Start() error {
	// TODO: Start the game with error checking based on G.state and G.completed
	if G.state == UNSTARTED {
		G.initialize()
	}
	G.state = RUNNING
	fmt.Println("Game running. (", G.player[WHITE].Name, "vs", G.player[BLACK].Name, ")") //TODO: include round#

	// Start up the engines:
	G.player[WHITE].Start()
	G.player[BLACK].Start()

	// Begin playing the game:
	G.playLoop()

	return nil
}

func (G *Game) Stop() error {
	// TODO: Stop the game with error checking based on G.state and G.completed
	if G.state == RUNNING {
		// Turn off the engines:
		G.player[WHITE].Shutdown()
		G.player[BLACK].Shutdown()

		G.state = STOPPED
		fmt.Println("Game stopped.")
	}
	return nil
}

/*******************************************************************************

	Functions that modify the game itself:

*******************************************************************************/

func (G *Game) MakeMove(m Move) error {
	// TODO: Make the move on G.board with proper error checking along the way

	G.moveList = append(G.moveList, m)

	// TODO: Actually change the location on the board =P

	return nil
}

func (G *Game) resetTimeControl() {
	G.movesToGo = G.moves
	G.timer = [2]int64{G.time, G.time}
}

func (G *Game) GameOver(looser Color) {
	G.result = []Color{BLACK, WHITE, DRAW}[looser] //opposite of the looser
	G.completed = true
	G.Stop()
}

func (G *Game) initialize() error {
	// Sets up the game so that its ready for white to make the first move

	// TODO: this assumes a fresh unstarted game.
	G.resetTimeControl()
	G.toMove = WHITE
	G.board.Reset()
	G.state = UNSTARTED
	return nil
}

func (G *Game) Validate() error {
	// TODO: Should check that all the data members are set up correctly to not cause a crash.
	return nil
}
