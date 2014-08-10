/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes, Daniel Sparks
 Created: 8/8/2014

 Module: Book
 Description: Opening Book. Right now it loads a pgn and plays the first few
 moves from the games in the PGN.

 TODO:
 	-adjust for carousel
 	-error handling
 	-PlayOpeniings() can go infinite

*******************************************************************************/

/*

KNOWN BUG:

Notation: Can not find source square.
r1bqkb1r/pppp1ppp/2n2n2/1B2p3/4P3/5N2/PPPP1PPP/RNBQK2R w KQkq - 2 4
d5
panic: runtime error: index out of range

*/

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"strings"
	"time"
)

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

// Play each game in the tourney enough to get out of the book:
func PlayOpenings(T *Tourney) error {

	//var FENs []string
	//var moves [][]Move
	//FENcount := []int{len(T.GameList) / 2, len(T.GameList)}[T.Rounds%2]
	//var BookIndexes []int

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
	/*
		for i := 0; i < FENcount; i++ {
			if T.RandomBook {
				r := rand.Intn(FENcount) % len(T.BookPGN)
				for onList(BookIndexes, r) {
					r = rand.Intn(FENcount)
				}
				BookIndexes = append(BookIndexes, r)
			} else {
				BookIndexes = append(BookIndexes, i%len(T.BookPGN))
			}
		}
	*/
	// Figure out what the FENs are from the pgn choices we made above:
	for i := 0; i < len(T.GameList); i++ {
		var dummy Game
		alreadyListed := true // hack to get the for loop to go at least once
		// Loop until a unique FEN is found:
		// TODO: This could go infinite if there arent enough unique positions.
		for alreadyListed {
			// pick a game from the book to use:
			var index int = i
			if T.RandomBook {
				index = rand.Intn(len(T.BookPGN))
			}
			// play out the opening on a dummy game:
			dummy = Game{}
			dummy.initialize()
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
		}
		// DEBUG:
		fmt.Println(dummy.FEN())

		// TODO : This will not work for for non-Carousel
		// Apply the FENs to the games on our game list:
		if !T.GameList[i].Completed {
			for _, mv := range dummy.MoveList {
				T.GameList[i].MakeMove(mv)
			}
			/*
				T.GameList[i].board = dummy.board
				T.GameList[i].enPassant = dummy.enPassant
				T.GameList[i].castleRights = dummy.castleRights
				T.GameList[i].fiftyRule = dummy.fiftyRule
				T.GameList[i].StartingFEN = dummy.FEN()
				T.GameList[i].MoveList = dummy.MoveList
			*/
		}
		if T.Rounds%2 == 0 {
			i++
			for _, mv := range dummy.MoveList {
				T.GameList[i].MakeMove(mv)
			}
		}
	}

	return nil
}
