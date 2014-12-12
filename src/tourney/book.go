/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes, Daniel Sparks
 Created: 8/8/2014

 Module: Book
 Description: Opening Book. Right now it loads a pgn and plays the first few
 moves from the games in the PGN.

 TODO:
 	- "Failed: Not enough unique positions" still allows the mirrored game to
 	  continue playing. should skip that one also.
 	-Count possible openings before playing
 	-adjust for carousel
 	-error handling
 	-Consolidate locations of error handling.

*******************************************************************************/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

const BOOKMOVE string = "Book Move"

func LoadBook(T *Tourney) error {
	// Check for pgn type:
	// TODO: check based on file contents not file name
	if strings.HasSuffix(T.BookLocation, ".pgn") {
		pgn, err := ioutil.ReadFile(T.BookLocation)
		if err != nil {
			return err
		}
		T.BookPGN = DecodePGN(string(pgn))
	}
	return nil
}
func CopyStartingPosition(From *Game, To *Game) error {
	To.StartingFEN = From.StartingFEN
	//Play the moves until the FEN matches the starting FEN:
	for i := 0; i < len(From.MoveList); i++ {
		if To.FEN() == To.StartingFEN {
			break
		}
		if err := To.MakeMove(From.MoveList[i]); err != nil {
			return err
		}

	}
	return nil
}

// Plays the opening from the pgn book for a single game:
func PlayOpening(T *Tourney, GameIndex int) error {

	// Helper function:
	alreadyUsed := func(n string) bool {
		for i, _ := range T.GameList {
			if n == T.GameList[i].StartingFEN {
				return true
			}
		}
		return false
	}

	rand.Seed(time.Now().Unix())
	if T.GameList[GameIndex].Completed {
		fmt.Println("Game already played.") //TODO: consolidate locaions of error handling.
		return nil
	}
	// Check if the opening has already been played on this game:
	if len(T.GameList[GameIndex].MoveList) > 0 {
		// TODO: better checking!
		fmt.Println("Opening already played.") //TODO: consolidate locaions of error handling.
		return nil
	}

	if T.Rounds%2 == 0 && GameIndex%2 == 1 && GameIndex > 0 {
		if T.BookMoves > 0 && T.GameList[GameIndex-1].StartingFEN == "" {
			return errors.New("No starting position to mirror.")
		}
		if err := CopyStartingPosition(&T.GameList[GameIndex-1], &T.GameList[GameIndex]); err != nil {
			return err
		}
		return nil
	}
	var dummy Game
	alreadyListed := true // hack to get the for loop to go at least once
	// Loop until a unique FEN is found:
	attempts := 0
	for alreadyListed {
		// escape if going infinite:
		if attempts > len(T.BookPGN) {
			return errors.New("Not enough unique positions in the opening book.")
		}
		// pick a game from the book to use:
		var index int = GameIndex
		if T.RandomBook {
			index = rand.Intn(len(T.BookPGN))
		}
		// play out the opening on a dummy game:
		dummy = T.GameList[GameIndex]
		for j := 0; j < 2*T.BookMoves; j++ {
			b := &T.BookPGN[index]
			// make sure we dont try to play more moves than what is in the book:
			if j >= len(b.MoveList) {
				break
			}
			mv := b.MoveList[j].Algebraic
			mv = StripAnnotations(mv)
			mv = InternalizeNotation(&dummy, mv)
			dummy.MakeMove(Move{Algebraic: mv, Comment: BOOKMOVE})
		}
		alreadyListed = alreadyUsed(dummy.FEN())
		attempts++
	}
	T.GameList[GameIndex] = dummy
	T.GameList[GameIndex].StartingFEN = dummy.FEN()
	/*
		if T.Rounds%2 == 0 && !T.GameList[GameIndex+1].Completed {
			dummy.Player = T.GameList[GameIndex+1].Player
			T.GameList[GameIndex+1] = dummy
			T.GameList[GameIndex+1].StartingFEN = dummy.FEN()
			GameIndex++
			// DEBUG:
			//fmt.Println(dummy.FEN())
		}
	*/
	// DEBUG:
	//fmt.Println(dummy.FEN())

	// TODO : This will not work for for non-Carousel

	return nil
}

/*
// Play each game in the tourney enough to get out of the book:
func PlayOpenings(T *Tourney) error {

	// Helper function:
	alreadyUsed := func(n string) bool {
		for i, _ := range T.GameList {
			if n == T.GameList[i].StartingFEN {
				return true
			}
		}
		return false
	}

	// Pick which games from the book pgn to use:
	rand.Seed(time.Now().Unix())

	//for i, _ := range T.GameList {
	for i := 0; i < len(T.GameList); i++ {
		if T.GameList[i].Completed {
			continue
		}
		var dummy Game
		alreadyListed := true // hack to get the for loop to go at least once
		// Loop until a unique FEN is found:
		attempts := 0
		for alreadyListed {
			// escape if going infinite:
			if attempts > len(T.BookPGN) {
				return errors.New("Not enough unique positions in the opening book.")
			}
			// pick a game from the book to use:
			var index int = i
			if T.RandomBook {
				index = rand.Intn(len(T.BookPGN))
			}
			// play out the opening on a dummy game:
			dummy = T.GameList[i]
			for j := 0; j < 2*T.BookMoves; j++ {
				b := &T.BookPGN[index]
				// make sure we dont try to play more moves than what is in the book:
				if j >= len(b.MoveList) {
					break
				}
				mv := b.MoveList[j].Algebraic
				mv = StripAnnotations(mv)
				mv = InternalizeNotation(&dummy, mv)
				dummy.MakeMove(Move{Algebraic: mv, log: []string{"Book Move."}})
			}
			alreadyListed = alreadyUsed(dummy.FEN())
			attempts++
		}
		T.GameList[i] = dummy
		T.GameList[i].StartingFEN = dummy.FEN()
		if T.Rounds%2 == 0 && !T.GameList[i+1].Completed {
			dummy.Player = T.GameList[i+1].Player
			T.GameList[i+1] = dummy
			T.GameList[i+1].StartingFEN = dummy.FEN()
			i++
			// DEBUG:
			//fmt.Println(dummy.FEN())
		}
		// DEBUG:
		//fmt.Println(dummy.FEN())

		// TODO : This will not work for for non-Carousel

	}

	return nil
}
*/
