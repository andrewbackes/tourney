/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes
 Created: 12/15/2014

 Module: Opening
 Description:

 TODO:

*******************************************************************************/

package main

import (
	"errors"
	//"fmt"
	//"io/ioutil"
	//"math/rand"
	//"strconv"
	//"strings"
	//"time"
)

func PlayOpening(T *Tourney, GameIndex int) error {

	if T.BookMoves <= 0 {
		return nil
	}

	if T.OpeningBook == nil {
		if book, err := LoadOrBuildBook(T.BookLocation, T.BookMoves, nil); err != nil {
			return err
		} else {
			T.OpeningBook = book
		}
	}

	/*
		if GameIndex == 0 {
			// need to build the Tourney's book iterator:
			if err := mapBookIterator(T, T.OpeningBook); err != nil {
				return err
			}
		}
	*/

	if T.GameList[GameIndex].Completed {
		return errors.New("Game already played.")
	}

	// Check if the opening has already been played on this game:
	if len(T.GameList[GameIndex].MoveList) > 0 {
		// TODO: better checking!
		//return errors.New("Opening already played.")
		return nil
	}

	var fen string
	var err error

	// Find the next opening to play:
	if fen, err = (T.OpeningBook).nextOpening(T, GameIndex); err != nil {
		return err
	}
	// Play it:
	opening := T.OpeningBook.Positions[T.BookMoves-1][fen]
	if err = applyOpeningToGame(opening, fen, &(T.GameList[GameIndex])); err != nil {
		return err
	}
	// All done!
	return nil
}

// Based on the tourney settings deduce what opening should be played
// for this matchup.
func (B *Book) nextOpening(T *Tourney, GameIndex int) (string, error) {
	// TODO: a lot.

	// Determine if we should just use the same position as the last game in the matchup:
	if mirror, fen := shouldMirror(T, GameIndex); mirror {
		return fen, nil
	}

	// find the last opening used and use the one after that:
	nextIndex := 0
	if GameIndex >= 1 {
		lastFEN := T.GameList[GameIndex-1].StartingFEN
		nextIndex = B.Positions[T.BookMoves-1][lastFEN].Index + 1
		if nextIndex >= len(B.Iterator[T.BookMoves-1]) {
			if T.RepeatOpenings {
				nextIndex = 0
			} else {
				return "", errors.New("Not enough openings in book.")
			}
		}
	}
	return B.Iterator[T.BookMoves-1][nextIndex], nil

	return "", nil
}

// Determines if the game should be mirrored and if so, returns what fen to use.
func shouldMirror(T *Tourney, GameIndex int) (bool, string) {

	if !T.BookMirroring || GameIndex < 1 {
		return false, ""
	}

	w := &(T.GameList[GameIndex].Player[WHITE])
	b := &(T.GameList[GameIndex].Player[BLACK])

	mirror := true
	fen := ""
	if GameIndex >= 1 {
		// is the previous game in this matchup
		lastgame := &(T.GameList[GameIndex-1])
		if lastgame.Player[WHITE].Equals(b) && lastgame.Player[BLACK].Equals(w) {
			// is it already a mirrored game?
			fen = lastgame.StartingFEN
			if GameIndex >= 2 {
				gameBeforeLast := &(T.GameList[GameIndex-2])
				if lastgame.StartingFEN == gameBeforeLast.StartingFEN {
					mirror = false
				}
			}
		} else {
			// last game is not in this matchup, so dont mirror it.
			mirror = false
		}
	}
	if !mirror {
		fen = ""
	}
	return mirror, fen
}

func applyOpeningToGame(opening BookPosition, fen string, G *Game) error {
	for _, move := range opening.MoveList {
		G.MakeMove(move)
		G.AddMoveAnalysis(MoveAnalysis{Comment: BOOKMOVE})
		
		if G.toMove() == WHITE {
			G.MovesToGo -= 1
			if G.MovesToGo == 0 && G.Repeating == true {
				G.resetTimeControl()
			}
		}
	}
	G.StartingFEN = fen
	return nil
}

/*
// Tourney.BookIteratorMap holds all of the indexes in Book.Iterator
// The idea is that when T.RandomBook is true the values of
// BookIteratorMap will be random values in Book.Iterator
// when RandomBook is false, it will just be a straight mapping of indexes.
func mapBookIterator(T *Tourney, B *Book) error {

	T.BookIteratorMap = B.Iterator[T.BookMoves]
	T.BookIteratorReverseMap = B.Iterator[T.BookMoves]

	if T.RandomBook {
		// TODO: randomize
	}

	return nil
}
*/
