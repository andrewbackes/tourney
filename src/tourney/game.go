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
	WHITE Color 	= 0
	BLACK Color 	= 1
	BOTH Color 		= 2
	DRAW Color 		= 2
)

type Game struct {

// Header info (General info, usually needed for pgn):
	time 			uint64
	moves 			uint64

// Pre game info:
	timer 			[2]int64 		//the game clock for each side. once =0, game is over
	movesRemaining 	[2]uint64 		//moves until time control
	player 			[2]Engine 		//white=0,black=1

// Game run time info:
	state 			State 			//RUNNING,STOPPED,PAUSED
	toMove 			Color 			//WHITE or BLACK only
	board 			Board
	moveList 		[]Move 			//move history

// Post game info:
	result 			Color 			//WHITE,BLACK,DRAW - should be set when game state is changed to STOPPED

}


//Methods:
//	start() - initializes needed game info (like settings the board, threading the engines, etc) and begins the loop
//	loop() - asks engine to move, updates game state, validates game state, repeats until end of game is reached
//	various utility methods, like movecount, etc.
//	makeMove()
//	load() - loads fen, pgn, etc
//	reset()

