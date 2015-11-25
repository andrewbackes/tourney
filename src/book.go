/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes
 Created: 8/8/2014

 Module: Book
 Description: Opening Book internal format. Loads a PGN file into a Book object.

 TODO:
 	-Count possible openings before playing
 	-adjust for carousel
 	-error handling
 	-Consolidate locations of error handling.

 Notes:
 	- Command used to create the default book:
 		buildbook /Users/Andrew/Documents/millionbase-2.22.pgn 14 Result =1/2-1/2 WhiteElo >2699 BlackElo >2699

*******************************************************************************/

package main

import (
	"errors"
	"fmt"
	//"io/ioutil"
	"math/rand"
	"strconv"
	"strings"
	//"time"
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

const BOOKMOVE string = "Book Move"

type BookPosition struct {
	Index    int
	Weight   int
	MoveList []Move
}

type Book struct {
	//Filename  string
	Positions []map[string]BookPosition
	Iterator  [][]string //keys
}

func NewBook(filename string, MaxMoves int) *Book {
	book := &Book{
		//Filename:  filename,
		Positions: make([]map[string]BookPosition, MaxMoves),
		Iterator:  make([][]string, MaxMoves),
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
	}
	B.sortByOccurrence()
}

func (B *Book) setIndexes() {
	for i, _ := range B.Iterator {
		for j, _ := range B.Iterator[i] {
			key := B.Iterator[i][j]
			value := B.Positions[i][key]
			value.Index = j
			B.Positions[i][key] = value
		}
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
func (B *Book) String() string {
	/*
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
					s += m + " "
				}
				s += fmt.Sprint(" (x", v.Weight, ") = [", i, "][", k, "]\n")
			}
			w += fmt.Sprint("weight([", i, "])=", weightSum, "; ")
		}

		return s + l + "\n" + w + "\n"
	*/
	var r string
	for d, _ := range B.Positions {
		r += "Move #" + strconv.Itoa(d+1) + ": " + strconv.Itoa(len(B.Positions[d])) + " positions.\n"
	}
	return r
}

// Tries to open the .book (json) version of the filename.
// When it cant be found it is built from the PGN .
//
// filename can be a .pgn or a .book
func LoadOrBuildBook(filename string, MoveNumber int, filters []PGNFilter) (*Book, error) {

	fmt.Print("Opening book: ", filename, "'...\n")
	// Look for the file in the given path
	if _, try := os.Stat(filename); try != nil {
		// when it cant be found in the given location, look in the books folder
		if _, try2 := os.Stat(filepath.Join(Settings.BookDirectory, filename)); try2 == nil {
			// file exists here, so use this path:
			filename = filepath.Join(Settings.BookDirectory, filename)
		} else {
			// cant find the file!
			return nil, errors.New("Can not find file: " + filename)
		}
	}

	bookfilename := ""
	if strings.HasSuffix(strings.ToLower(filename), ".pgn") {
		// PGN file, so we need to match a .book file to it:
		bookfilename = filename[:len(filename)-len(".pgn")] + ".book"
	} else if strings.HasSuffix(strings.ToLower(filename), ".book") {
		// already given a .book file:
		bookfilename = filename
	} else {
		return nil, errors.New("Invalid book format: '" + filename + "' does not end with .pgn or .book")
	}

	// Now we look for the .book file:
	foundit := false
	fmt.Print("Looking for previously build book '" + bookfilename + "'...\n")
	if _, try := os.Stat(bookfilename); try != nil {
		// when it cant be found in the given location, look in the books folder
		_, f := filepath.Split(bookfilename)
		bookfilename = filepath.Join(Settings.BookDirectory, f)
		fmt.Print("Looking for previously build book '", bookfilename, "'...\n")
		if _, try2 := os.Stat(bookfilename); try2 == nil {
			//file exists here, so use this path:
			foundit = true
		}
	} else {
		foundit = true
	}
	var b *Book
	var e error

	if foundit {
		// when we find the already built book, we just need to load it:
		//fmt.Println("Found it.")
		b, e = OpenBook(bookfilename)
	} else {
		// couldn't find the .book, so we need to build it:
		b, e = BuildBookFromPGN(filename, MoveNumber, filters)
		if e == nil {
			if err := b.SaveBook(bookfilename); err != nil {
				fmt.Println(err)
			}
		}
	}
	return b, e
}

// Opens a .book file:
func OpenBook(filename string) (*Book, error) {
	// Try to open the file:
	fmt.Print("Using opening book: '", filename, "'...\n")
	bookFile, err := os.Open(filename)
	defer bookFile.Close()
	if err != nil {
		fmt.Print("Failed to open: '", filename, "'\n", err.Error(), "\n")
		return nil, err
	}
	// Make the object:
	book := NewBook(filename, 0) // TODO: should not be fixed!
	// Try to decode the file:
	jsonParser := json.NewDecoder(bookFile)
	if err = jsonParser.Decode(book); err != nil {
		fmt.Println("Failed to decode:", err.Error())
		return nil, err
	}
	return book, nil
}

func (B *Book) SaveBook(filename string) error {

	// Create the Save directory:
	if Settings.BookDirectory != "" {
		if err := os.MkdirAll(Settings.BookDirectory, os.ModePerm); err != nil {
			fmt.Println("Could not make directory:", Settings.BookDirectory, "\n", err)
			return err
		}
	}

	//check if the file exists:

	//filename = filepath.Join(Settings.BookDirectory, filename)

	fmt.Println("Saving '" + filename + "'... ")
	var file *os.File
	var err error
	if _, er := os.Stat(filename); os.IsNotExist(er) {
		// file doesnt exist
	} else if er == nil {
		// file does exist
		os.Remove(filename)
	}

	file, err = os.Create(filename)
	defer file.Close()

	var encoded []byte
	encoded, err = json.Marshal(*B)
	if err != nil {
		return err
	}
	if _, err = file.Write(encoded); err != nil {
		return err
	}

	return nil
}

// Load the PGN file into a Book object:
func BuildBookFromPGN(PGNfilename string, MoveNumber int, filters []PGNFilter) (*Book, error) {

	book := NewBook(PGNfilename, MoveNumber)

	// load the pgn games:
	PGN, err := ReadPGN(PGNfilename, filters) //LoadPGN(PGNfilename)
	if err != nil {
		return nil, err
	}

	// Progress bar:
	fmt.Println("Building Opening Book from " + PGNfilename + "...")
	fmt.Print("1%", strings.Repeat(" ", 36), "50%", strings.Repeat(" ", 35), "100%\n")
	dotgap := (len(*PGN) / 80)
	if len(*PGN)%80 != 0 {
		dotgap++
	}

	// go through each game in the pgn. get the fen at each move.
	// save the movelist
	for i, _ := range *PGN {
		dummyGame := NewGame()
		for ply := 1; ply <= 2*MoveNumber; ply += 2 {
			//for ply := 1; ply <= len((*PGN)[i].MoveList); ply += 2 {
			//if ply > 2*MoveNumber {
			//	break
			//}
			//check if this game has enough moves made:
			if len((*PGN)[i].MoveList) < ply+1 {
				//fmt.Println("len((*PGN)[i].MoveList) < ply+1")
				break
			}
			//white move:
			wmv, err := ConvertToPCN(&dummyGame, StripAnnotations(string((*PGN)[i].MoveList[ply-1])))
			if err != nil {
				break
			}
			if err := dummyGame.MakeMove(Move(wmv)); err != nil {
				break
			}
			//black move:
			bmv, err := ConvertToPCN(&dummyGame, StripAnnotations(string((*PGN)[i].MoveList[ply])))
			if err != nil {
				break
			}
			if err := dummyGame.MakeMove(Move(bmv)); err != nil {
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
		// update the console:
		if (i % dotgap) == 0 {
			fmt.Print(".")
		}
	}
	book.Iterate()
	fmt.Print("\n")
	return book, nil
}

func (B *Book) uniquePositions(BookMoves int) int {
	return len(B.Positions[BookMoves-1])
}

func (B *Book) Randomize(seed int64) {
	r := rand.New(rand.NewSource(seed))
	swap := func(depth, i, j int) {
		key := B.Iterator[depth][i]
		B.Iterator[depth][i] = B.Iterator[depth][j]
		B.Iterator[depth][j] = key
	}
	for depth, _ := range B.Iterator {
		for i := len(B.Iterator[depth]) - 1; i>=0; i-- {
			j := r.Intn(i+1)
			swap(depth, i, j)
		}
	}
	B.setIndexes()
}

type sortWrapper struct {
	book *Book
	depth int
}

func (sw sortWrapper) Len() int {
	return len(sw.book.Iterator[sw.depth])
}

func (sw sortWrapper) Swap(i, j int) {
	key := sw.book.Iterator[sw.depth][i]
	sw.book.Iterator[sw.depth][i] = sw.book.Iterator[sw.depth][j]
	sw.book.Iterator[sw.depth][j] = key
}

func (sw sortWrapper) Less(i, j int) bool {
	weightOf := func(i int)int {
		key := sw.book.Iterator[sw.depth][i]
		value := sw.book.Positions[sw.depth][key]
		return value.Weight
	} 
	return weightOf(i) < weightOf(j)
}

func (B *Book) sortByOccurrence() {
	for depth, _ := range B.Iterator {
		sw := sortWrapper{
			book: B,
			depth: depth,
		}
		sort.Sort(sw)
	}
	B.setIndexes()
}