/*

 Project: Tourney

 Module: games
 Description: holds the game object and methods for interacting with it

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*/

package main

type Color uint8

const (
	WHITE Color = 0
	BLACK Color = 1
	BOTH  Color = 2
	DRAW  Color = 2
)

type Game struct {

	// Header info (General info, usually needed for pgn):
	time  uint64
	moves uint64

	// Pre game info:
	timer          [2]int64  //the game clock for each side. once =0, game is over
	movesRemaining [2]uint64 //moves until time control
	player         [2]Engine //white=0,black=1

	// Game run time info:
	state     Status //RUNNING,STOPPED,PAUSED
	toMove    Color  //WHITE or BLACK only
	board     Board
	moveList  []Move //move history
	completed bool

	// Post game info:
	result Color //WHITE,BLACK,DRAW - should be set when game state is changed to STOPPED

}

//Methods:
//	loop() - asks engine to move, updates game state, validates game state, repeats until end of game is reached
//	MAYBE: various utility methods, like movecount, etc.
//	MAYBE: load() - loads fen, pgn, etc
//	reset()

func (G *Game) Validate() error {
	// TODO: Should check that all the data members are set up correctly to not cause a crash.
	return nil
}

func (G *Game) Start() error {
	// TODO: Start the game with error checking based on G.state and G.completed
	return nil
}

func (G *Game) Stop() error {
	// TODO: Stop the game with error checking based on G.state and G.completed

	return nil
}

func (G *Game) MakeMove(m Move) error {
	// TODO: Make the move on G.board with proper error checking along the way

	return nil
}
