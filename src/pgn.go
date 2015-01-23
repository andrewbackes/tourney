/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes
 Created: 8/8/2014

 Module: PGN
 Description: PGN tools

 TODO:
 		-finish tags: ELO, time, timecontrol

*******************************************************************************/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

type PGNGame struct {
	Tags     map[string]string
	MoveList []Move
}

type PGNFilter struct {
	Tag   string // ex: WhiteElo
	Value string // ex: ">2700"
}

func NewPGNGame() PGNGame {
	return PGNGame{
		Tags: make(map[string]string),
	}
}

//Load the contents of a PGN file into memory:
func LoadPGN(filename string) (*[]Game, error) {
	if !strings.HasSuffix(filename, ".pgn") {
		return nil, errors.New("Invalid PGN file.")
	}
	pgn, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	PGN := DecodePGN(string(pgn))
	return &PGN, nil
}

// Turns Game structs into PGN
func EncodePGN(G *Game) string {
	// TODO: Test needed. Changed code about Move.log without testing. See below.

	var pgn string
	tags := [][]string{
		{"Event", G.Event},
		{"Site", G.Site},
		{"Date", G.Date},
		{"Round", strconv.Itoa(G.Round)},
		{"White", G.Player[WHITE].Name},
		{"Black", G.Player[BLACK].Name},
		{"Result", "*"},
		{"WhiteElo", "-"},
		{"BlackElo", "-"},
		{"Time", "-"},
		{"TimeControl", "-"},
	}
	if G.StartingFEN != "" {
		tags = append(tags, []string{"Setup", "1"})
		tags = append(tags, []string{"FEN", G.StartingFEN})
	}
	if G.Completed {
		tags[6][1] = []string{"1-0", "0-1", "1/2-1/2"}[G.Result]
	}
	for _, t := range tags {
		pgn += fmt.Sprintln("[" + t[0] + " \"" + t[1] + "\"]")
	}
	pgn += fmt.Sprintln("")

	for j, _ := range G.MoveList {
		// TODO: replaced this code without testing:
		//if len(G.MoveList[j].log) > 0 && strings.Contains(G.MoveList[j].log[0], "Book Move.") {
		if G.MoveList[j].Comment == BOOKMOVE {
			// dont print book moves, since the FEN tag would mess it up.
			continue
		}
		if j%2 == 0 {
			pgn += strconv.Itoa((j/2)+1) + ". "
		}
		pgn += G.MoveList[j].Algebraic + " "
	}
	pgn += tags[6][1]
	pgn += fmt.Sprintln("")
	pgn += fmt.Sprintln("")

	return pgn
}

// Turns PGN into Game Structs
func DecodePGN(pgn string) []Game {
	var gamelist []Game

	game := strings.Split(pgn, "[Event")
	for i, _ := range game {
		var G Game
		if game[i] != "" {
			game[i] = "[Event" + game[i]
		}
		lines := strings.Split(game[i], "\n")
		for _, l := range lines {
			if l == "" {
				continue
			}
			if strings.HasPrefix(l, "[") {
				// its a tag
				t := strings.SplitN(l, " ", 2)
				key := t[0][1:]
				value := t[1][1 : len(t[1])-2] // TODO: index out of range error here.
				switch key {
				case "Event":
					G.Event = value
				case "Site":
					G.Site = value
				case "Date":
					G.Date = value
				case "Round":
					G.Round, _ = strconv.Atoi(value)
				case "White":
					G.Player[WHITE].Name = value
				case "Black":
					G.Player[BLACK].Name = value
				case "Result":
					G.Result = map[string]Color{"1-0": WHITE, "0-1": BLACK, "1/2-1/2": DRAW, "*": NEITHER}[value]
				}
			} else {
				// its a move
				moves := strings.Split(l, " ")
				for _, m := range moves {
					if strings.HasPrefix(m, "{") {
						continue
					}
					if strings.Contains(m, ".") {
						m = strings.Split(m, ".")[1]
						m = strings.Trim(m, " ")
					}
					//SAN := "([BKNPQR]?)([a-h]?)([0-9]?)([x=]?)([BKNPQR]|[a-h][1-8])([+#!?]?)([+#!?]?)"
					//if ok, _ := regexp.MatchString(SAN, m); ok || m == "O-O" || m == "O-O-O" {

					if isMove(StripAnnotations(m)) {
						G.MoveList = append(G.MoveList, Move{Algebraic: m})
					}
				}
			}
		}
		if (G.Event != "") && (len(G.MoveList) > 0) {
			gamelist = append(gamelist, G)
		}
	}

	return gamelist
}

// Line by line reads a pgn file. Turns what is read into a PGNGame struct.
// This is an improvement over LoadPGN() since it goes line by line.
// For huge PGN files, this will work but LoadPGN() will not.
func ReadPGN(filename string, filters []PGNFilter) (*[]PGNGame, error) {

	fmt.Println("Reading", filename, "...")
	if !strings.HasSuffix(filename, ".pgn") {
		return nil, errors.New("Invalid PGN file.")
	}

	// Open the file:
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	GameList := []PGNGame{}

	// Display progress bar:
	var filesize int64
	if fs, err := file.Stat(); err == nil {
		filesize = fs.Size()
	}
	dotsPrinted := 0
	bytesPerDot := filesize / 80
	fmt.Print("1%", strings.Repeat(" ", 36), "50%", strings.Repeat(" ", 35), "100%\n")
	updateProgressBar := func(completed int64) {
		dotnumber := int(completed / bytesPerDot)
		if (dotsPrinted) < dotnumber {
			dotsPrinted++
			fmt.Print(".")
		}
	}

	// Read line by line:
	var bytesRead int64
	scanner := bufio.NewScanner(file)
	readingmoves := false // flag
	currentGame := NewPGNGame()
	for scanner.Scan() {
		line := scanner.Bytes()
		bytesRead += int64(len(line))
		if len(line) == 0 {
			continue
		}
		if line[0] == '[' {
			if readingmoves {
				// since we are no longer reading moves, we know this is a new game
				readingmoves = false
				// so we need to sort out what to do with the game that we previously read:
				if SatisfiesFilters(&currentGame, &filters) {
					GameList = append(GameList, currentGame)
				}
				currentGame = NewPGNGame()
				updateProgressBar(bytesRead)
			}
			key, value := SplitTag(line)
			currentGame.Tags[string(key)] = string(value)
		} else {
			readingmoves = true
			appendMoves(&currentGame, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	fmt.Println("\nRead", len(GameList), "games in PGN.")
	return &GameList, nil
}

func appendMoves(game *PGNGame, line []byte) {
	// example: 1. e2e4 d7d5 2. b1c3 f7f5 {asd asd} 3. a2a3 ;asdasdasdasd"
	l := RemoveComments(line)
	//l = RemoveNumbering(line)
	moves := strings.Split(string(l), " ")
	for _, m := range moves {
		if strings.Contains(m, ".") {
			m = strings.Split(m, ".")[1]
			m = strings.Trim(m, " ")
		}
		//if isMove(StripAnnotations(m)) {
		if m != "1/2-1/2" && m != "1-0" && m != "0-1" && m != "*" && m != "" {
			game.MoveList = append(game.MoveList, Move{Algebraic: m})
		}
		//}

	}
}

func RemoveComments(line []byte) []byte {
	marker := 0
	i := 0
	for {
		if line[i] == ';' {
			if i-1 < 0 {
				return []byte{}
			}
			return line[:i-1]
		}
		if line[i] == '{' {
			marker = i
		}
		if line[i] == '}' {
			if marker-1 < 0 {
				marker = marker + 1
			}
			if i+1 > len(line)-1 {
				line = line[:marker-1]
			} else {
				line = append(line[:marker-1], line[i+1:]...)
			}
			i = marker
		}
		i++
		if i > len(line)-1 {
			break
		}
	}
	return line
}

func RemoveNumbering(line []byte) []byte {
	marker := 0
	i := 0
	for {
		if line[i] == ' ' {
			marker = i
		}
		if line[i] == '.' {
			if i+1 > len(line)-1 {
				line = line[:marker]
			} else if marker == 0 && i+2 < len(line) {
				line = line[i+2:]
			} else {
				line = append(line[:marker], line[i+1:]...)
			}
			i = marker
		}

		i++
		if i > len(line)-1 {
			break
		}
	}
	return line
}

// takes a tag and returns its key and value components.
// ex: [Event "Testing"] ==> "Event", "Testing"
func SplitTag(line []byte) ([]byte, []byte) {

	// Look for a space:
	space := 0
	for i, _ := range line {
		if line[i] == ' ' {
			space = i
			break
		}
	}
	key := line[1:space]
	value := line[space+2 : len(line)-2]

	return key, value
}

// Verifies that the game meets the requirements of the filter.
func SatisfiesFilters(game *PGNGame, filters *[]PGNFilter) bool {
	if filters == nil {
		return true
	}
	for _, f := range *filters {
		v := game.Tags[f.Tag]
		if v == "" || !ValuesMatch(v, f.Value) {
			return false
		}
	}
	return true
}

// ex: ValuesMatch("2701",">2700") == true
func ValuesMatch(value string, constraint string) bool {
	if len(constraint) > 1 {
		c := constraint[1:]
		switch constraint[:1] {
		case ">":
			c_int, _ := strconv.Atoi(c)
			value_int, _ := strconv.Atoi(value)
			return value_int > c_int
		case "<":
			c_int, _ := strconv.Atoi(c)
			value_int, _ := strconv.Atoi(value)
			return value_int < c_int
		case "=":
			return value == c
		}
	}
	return value == constraint
}
