/*

 Project: Tourney

 Module: Moves
 Description: holds the move object and methods for interacting with it.
 	Eventually, engine data/logs will be tied into this?

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*/

package main

type Move struct {
	Algebraic string
	log       []string
}

func getMove(from uint, to uint) Move {
	// makes a move object from the to/from square index
	var r Move
	r.Algebraic = getAlg(from) + getAlg(to)
	return r
}
