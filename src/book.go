/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes
 Created: 8/8/2014

 Module: Book
 Description: Opening Book. Loads a PGN file into a Book object.

 TODO:
 	-Count possible openings before playing
 	-adjust for carousel
 	-error handling
 	-Consolidate locations of error handling.

*******************************************************************************/

package main

import (
	//"errors"
	"fmt"
	//"io/ioutil"
	//"math/rand"
	"strconv"
	//"strings"
	//"time"
)

const BOOKMOVE string = "Book Move"

type BookPosition struct {
	Index    int
	Weight   int
	MoveList []Move
}

type Book struct {
	FromFilename string
	Positions    []map[string]BookPosition
	Iterator     [][]string //keys
}

func NewBook(PGNFilename string, MaxMoves int) *Book {
	book := &Book{
		FromFilename: PGNFilename,
		Positions:    make([]map[string]BookPosition, MaxMoves),
		Iterator:     make([][]string, MaxMoves),
	}
	for m := 0; m < MaxMoves; m++ {
		book.Positions[m] = make(map[string]BookPosition)
	}
	return book
}

func (B *Book) Iterate() {
	for i, _ := range B.Positions {
		for k, _ := range B.Positions[i] {
			B.Iterator[i] = append(B.Iterator[i], k)
			value := B.Positions[i][k]
			value.Index = len(B.Iterator[i]) - 1
			B.Positions[i][k] = value
		}
		// TODO: sort by weight. then reassign the index's
	}
}

func NewBookPosition(Moves []Move) *BookPosition {
	return &BookPosition{
		Index:    0,
		Weight:   1,
		MoveList: Moves,
	}
}

// For use in the fmt package. Tell's Print how to display the Book
// object.
func (B Book) String() string {
	var s, l, w string
	s += fmt.Sprint("Moves#\tFEN\n")
	for i, _ := range B.Iterator {
		weightSum := 0
		l += fmt.Sprint("len([", i, "])=", len(B.Iterator[i]), "; ")
		//for k, v := range B.Positions[i] {
		for _, k := range B.Iterator[i] {
			v := B.Positions[i][k]
			weightSum += v.Weight
			s += strconv.Itoa(v.Index) + ".\t"
			for _, m := range v.MoveList {
				s += m.Algebraic + " "
			}
			s += fmt.Sprint(" (x", v.Weight, ") = [", i, "][", k, "]\n")
		}
		w += fmt.Sprint("weight([", i, "])=", weightSum, "; ")
	}

	return s + l + "\n" + w + "\n"
}

// Tries to load the json version of the book that would have been built
// by this pgn file's name. When it can't find it or it is invalid
// somehow, then it builds a new json verion based on the pgn.
func LoadOrBuildBook(PGNfilename string, MoveNumber int) (*Book, error) {
	// TODO: complete this function.

	// try to load internal format of book:

	// if not build it:
	return BuildBook(PGNfilename, MoveNumber)

	// save what was built.

}

// Load the PGN file into a Book object:
func BuildBook(PGNfilename string, MoveNumber int) (*Book, error) {
	fmt.Print("Building Opening Book...")

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
			wmv, err := ConvertToPCN(&dummyGame, StripAnnotations((*PGN)[i].MoveList[ply-1].Algebraic))
			if err != nil {
				break
			}
			if err := dummyGame.MakeMove(Move{Algebraic: wmv}); err != nil {
				break
			}
			//black move:
			bmv, err := ConvertToPCN(&dummyGame, StripAnnotations((*PGN)[i].MoveList[ply].Algebraic))
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
	}
	book.Iterate()

	return book, nil
}

func (B *Book) uniquePositions(BookMoves int) int {
	return len(B.Positions[BookMoves-1])
}

/*
// figure out if there are enough unique positions in the opening book
// for each matchup to get a different one.
func EnoughPositionsInBook(T *Tourney, B *Book) bool {
	// TODO: This function does not work as intended!!!
	// 		 ex: 3 rounds, repeating, mirroring.

	if T.RepeatOpenings {
		roundsAsWhite := T.Rounds / 2
		if T.Rounds%2 == 1 {
			roundsAsWhite = (T.Rounds + 1) / 2
		}
		if T.BookMirroring {
			return (B.uniquePositions(T.BookMoves) >= roundsAsWhite)
		}
		return (B.uniquePositions(T.BookMoves) >= T.Rounds)
	}
	maxUniqueRounds := len(T.GameList)
	if T.BookMirroring && T.Rounds > 1 {
		maxUniqueRounds += (maxUniqueRounds % 2)
		maxUniqueRounds = maxUniqueRounds / 2
	}
	return B.uniquePositions(T.BookMoves) >= maxUniqueRounds
}
*/

// **********************
// OLD CODE:
// **********************

/*
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
*/

/*
// Plays the opening from the pgn book for a single game:
func PlayOpeningOLD(T *Tourney, GameIndex int) error {

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
			mv, err = ConvertToPCN(&dummy, mv)
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

	// TODO : This will not work for for non-Carousel

	return nil
}
*/
