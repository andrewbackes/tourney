/*

 Project: Tourney

 Module: Moves
 Description: holds the move object and methods for interacting with it.
 	Eventually, engine data/logs will be tied into this?

 Author(s): Andrew Backes
 Created: 7/16/2014

*/

package main

const MATESCORE int = 100000

type Move struct {
	Algebraic  string
	Ponder     string           `json:",omitempty"`
	Comment    string           `json:",omitempty"`
	Evaluation []EvaluationData `json:",omitempty"`
}

func (M *Move) Depth() int {
	if len(M.Evaluation) > 0 {
		return M.Evaluation[len(M.Evaluation)-1].Depth
	}
	return 0
}

func (M *Move) Pv() string {
	if len(M.Evaluation) > 0 {
		// look for the last PV that was stored:
		for i := len(M.Evaluation) - 1; i >= 0; i-- {
			if M.Evaluation[i].Pv != "" {
				return M.Evaluation[i].Pv
			}
		}
	}
	return ""
}

func (M *Move) Score() int {
	if len(M.Evaluation) > 0 {
		return M.Evaluation[len(M.Evaluation)-1].Score
	}
	return 0
}

func (M *Move) Time() int {
	if len(M.Evaluation) > 0 {
		return M.Evaluation[len(M.Evaluation)-1].Time
	}
	return 0
}

func (M *Move) Nodes() int {
	if len(M.Evaluation) > 0 {
		return M.Evaluation[len(M.Evaluation)-1].Nodes
	}
	return 0
}

// *******************************************************************

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
