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

// **********************
// NEW CODE:
// **********************

type BookPosition struct {
	Weight   int
	MoveList []Move
}

type Book struct {
	FromFilename string
	Positions    []map[string]BookPosition
}

func NewBook(PGNFilename string, MaxMoves int) *Book {
	book := &Book{
		FromFilename: PGNFilename,
		Positions:    make([]map[string]BookPosition, MaxMoves),
	}
	for m := 0; m < MaxMoves; m++ {
		book.Positions[m] = make(map[string]BookPosition)
	}
	return book
}

func NewBookPosition(Moves []Move) *BookPosition {
	return &BookPosition{
		Weight:   1,
		MoveList: Moves,
	}
}

// for use in the fmt package:
func (B Book) String() string {
	var s, l, w string
	s += fmt.Sprint("Moves#\tFEN\n")
	for i, _ := range B.Positions {
		sum := 0
		l += fmt.Sprint("len([", i, "])=", len(B.Positions[i]), "; ")
		for k, v := range B.Positions[i] {
			sum += v.Weight
			for _, m := range v.MoveList {
				s += m.Algebraic + " "
			}
			s += fmt.Sprint(" (x", v.Weight, ") = [", i, "][", k, "]\n")
		}
		w += fmt.Sprint("weight([", i, "])=", sum, "; ")
	}

	return s + l + "\n" + w + "\n"
}

func BuildBook(PGNfilename string, MoveNumber int) (*Book, error) {
	book := NewBook(PGNfilename, MoveNumber)

	// load the pgn games:
	PGN, err := LoadPGN(PGNfilename)
	if err != nil {
		return nil, err
	}

	// go through each game in the pgn. get the fen at each move.
	// save the movelist

	for i, _ := range *PGN {
		dummyGame := NewGame()
		for ply := 1; ply <= 2*MoveNumber; ply += 2 {
			//check if this game has enough moves made:
			if len((*PGN)[i].MoveList) < ply+1 {
				break
			}
			//white move:
			wmv, err := InternalizeNotation(&dummyGame, StripAnnotations((*PGN)[i].MoveList[ply-1].Algebraic))
			if err != nil {
				break
			}
			if err := dummyGame.MakeMove(Move{Algebraic: wmv}); err != nil {
				break
			}
			//black move:
			bmv, err := InternalizeNotation(&dummyGame, StripAnnotations((*PGN)[i].MoveList[ply].Algebraic))
			if err != nil {
				break
			}
			if err := dummyGame.MakeMove(Move{Algebraic: bmv}); err != nil {
				break
			}
			//get fen:
			fen := dummyGame.FEN()
			//add it to our book:
			if value, exists := book.Positions[ply/2][fen]; exists {
				value.Weight++
				book.Positions[ply/2][fen] = value
			} else {
				book.Positions[ply/2][fen] = *NewBookPosition(dummyGame.MoveList)
			}
		}
		//break //temporary
	}

	return book, nil
}

// **********************
// OLD CODE:
// **********************

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
			var err error
			mv, err = InternalizeNotation(&dummy, mv)
			if err != nil {
				return errors.New("Book notation error: '" + err.Error())
			}
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
