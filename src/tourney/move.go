/*

 Project: Tourney

 Module: Moves
 Description: holds the move object and methods for interacting with it.
 	Eventually, engine data/logs will be tied into this?

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*/

package main

const MATESCORE int = 100000

// TODO: Mate scores should be indicated as 100000 + N for "mate in N moves", and -100000 - N for "mated in N moves".

type Move struct {
	Algebraic string
	Ponder    string
	Comment   string
	//Depth     int
	//Score     int
	//Time      int
	//Nodes     int
	//Pv        string
	Evaluation []EvaluationData
}

func getMove(from uint, to uint) Move {
	// makes a move object from the to/from square index
	var r Move
	r.Algebraic = getAlg(from) + getAlg(to)
	return r
}

func MateIn(MovesTillMate int) int {
	if MovesTillMate < 0 {
		return MovesTillMate - MATESCORE
	}
	return MovesTillMate + MATESCORE
}
