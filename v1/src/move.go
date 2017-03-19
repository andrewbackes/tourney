/*

 Project: Tourney

 Module: Moves
 Description: holds the move object and methods for interacting with it.
 	Eventually, engine data/logs will be tied into this?

 Author(s): Andrew Backes
 Created: 7/16/2014

*/

package main

import (
	"strings"
)

const MATESCORE int = 100000

type Move string

type MoveAnalysis struct {
	Ponder     string           `json:",omitempty"`
	Comment    string           `json:",omitempty"`
	Evaluation []EvaluationData `json:",omitempty"`
}

func (M MoveAnalysis) Depth() int {
	if len(M.Evaluation) > 0 {
		return M.Evaluation[len(M.Evaluation)-1].Depth()
	}
	return 0
}

func (M MoveAnalysis) Pv() string {
	if len(M.Evaluation) > 0 {
		// look for the last PV that was stored:
		for i := len(M.Evaluation) - 1; i >= 0; i-- {
			if M.Evaluation[i].Pv() != "" {
				return M.Evaluation[i].Pv()
			}
		}
	}
	return ""
}

// PvChanges counts the number of times the PV changed for a move.
func (M MoveAnalysis) PvChanges() int {
	count := 0
	for i := 1; i < len(M.Evaluation); i++ {
		if !strings.Contains(strings.TrimSpace(M.Evaluation[i].Pv()), strings.TrimSpace(M.Evaluation[i-1].Pv())) {
			count++
		}
	}
	return count
}

// tallyStat is a helper function to count the occurances of a stat in the
// evaluation data.
func (M MoveAnalysis) tallyStat(hasStat func(int) bool) int {
	count := 0
	for i, _ := range M.Evaluation {
		if hasStat(i) {
			count++
		}
	}
	return count
}

// Lowerbounds counts the number of times a lowerbound was declared in
// the evaluation data.
func (M MoveAnalysis) Lowerbounds() int {
	return M.tallyStat(func(i int) bool { return M.Evaluation[i].Lowerbound() })
}

// Upperbounds counts the number of times a Upperbound was declared in
// the evaluation data.
func (M MoveAnalysis) Upperbounds() int {
	return M.tallyStat(func(i int) bool { return M.Evaluation[i].Upperbound() })
}

func (M MoveAnalysis) Score() int {
	if len(M.Evaluation) > 0 {
		return M.Evaluation[len(M.Evaluation)-1].Score()
	}
	return 0
}

func (M MoveAnalysis) Time() int {
	if len(M.Evaluation) > 0 {
		return M.Evaluation[len(M.Evaluation)-1].Time()
	}
	return 0
}

func (M MoveAnalysis) Nodes() int {
	if len(M.Evaluation) > 0 {
		return M.Evaluation[len(M.Evaluation)-1].Nodes()
	}
	return 0
}

// *******************************************************************

func getMove(from uint, to uint) Move {
	// makes a move object from the to/from square index
	var r Move
	r = Move(getAlg(from) + getAlg(to))
	return r
}

func MateIn(MovesTillMate int) int {
	if MovesTillMate < 0 {
		return MovesTillMate - MATESCORE
	}
	return MovesTillMate + MATESCORE
}
