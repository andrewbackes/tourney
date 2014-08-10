/*******************************************************************************

 Project: Tourney

 Module: tourney

 Description: A Tourney object impliments the Game object. The Engine object is
 also implimented but is used for its data fields, its methods are ignored. The
 tournament takes place in the playLoop() method. The Start() and Stop() methods
 essentially modify the data feild "state" which is read by playLoop().

 TODO:
 	-Allow for games to be distributed to multiple machines to be played.
 		-Each machine will have to be benchmarked to determine equivalent
 		 time control parameters.
 	-More tournament parameters

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*******************************************************************************/

package main

import (
	//"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strings"
	//"time"
	//"runtime"
)

type Status int

const (
	UNSTARTED Status = iota
	RUNNING          // in progress
	STOPPED
)

type Tourney struct {
	// Predetermined Settings for a tourney:
	Event string //identifier for this tournament. Unique may be better?
	Site  string
	Date  string

	Engines []Engine // which engines are playing in the tournament

	// The following will determine gauntlet, multigauntlet, roundrobin
	// 		if testSeats=1 then normal gauntlet (for the first engine)
	// 		if testSeats=#Engines then its roundrobine
	// 		if testSeats=2 then the first 2 engines will be multigauntlet
	TestSeats int

	Carousel bool //The order the engines play against eachother

	// Time control (Moves, time, repeating):
	Moves     int64 // moves per time control
	Time      int64 // time per time control in milliseconds
	Repeating bool  // restart time after moves hits

	Rounds int //number of games each engine will play

	// Opening book information:
	BookLocation string // File location of the book
	BookMoves    int    // Number of moves to use out of the book
	BookPGN      []Game
	RandomBook   bool

	QuitAfter bool //Quit after the tourney is complete.

	// Control settings (Determined while tourney is running, or when the tourney starts)
	//State     Status //flag to indicate: running, paused, stopped
	StateFlow chan Status
	GameList  []Game //list of all games in the tourney. populated when the tourney starts
	//activeGame *Game  //points to the currently running game in the list. Rethink this for multiple running games at a later time.

}

type Record struct {
	Player     Engine
	Wins       int
	Losses     int
	Draws      int
	Incomplete int
}

func RunTourney(T *Tourney) error {
	var state Status
	for i, _ := range T.GameList {
		select {
		case state = <-T.StateFlow:
			switch state {
			case RUNNING:
				fmt.Println("Tourney is running.")
			case STOPPED:
				fmt.Println("Tourney is stopped.")
			default:
				return errors.New("Tourney was put into an unknown state.")
			}
		default:
			// Dont block!
		}
		if state == STOPPED {
			// More clean up code here.
			break
		}
		// Start the next game on the list:
		fmt.Println("Round ", i, ": ", T.GameList[i].Player[WHITE].Name, "vs", T.GameList[i].Player[BLACK].Name)
		if !T.GameList[i].Completed {
			if err := PlayGame(&T.GameList[i]); err != nil {
				fmt.Println(err.Error())
				break
			}
			// DEBUG:
			//for _, mv := range T.GameList[i].MoveList {
			//	fmt.Println(mv.Algebraic)
			//}
			T.GameList[i].PrintHUD()
		}

	}
	// Empty the channel: (TODO: this may not be needed anymore)
	emptied := false
	for !emptied {
		select {
		case <-T.StateFlow:
		default:
			emptied = true
		}
	}

	// Show results:
	ShowResults(T)
	// Save results
	if err := SaveResults(T); err != nil {
		fmt.Println(err)
	}
	return nil
}

func SaveResults(T *Tourney) error {
	//check if the file exists:
	filename := T.Event + ".results"
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
	encoded, err = json.MarshalIndent(T.GameList, "", "  ")
	if err != nil {
		return err
	}
	if _, err = file.Write(encoded); err != nil {
		return err
	}
	fmt.Println("Successfully saved " + filename)
	return nil
}

func LoadPreviousResults(T *Tourney) (bool, error) {
	filename := T.Event + ".results"
	var err error
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// file doesnt exist
		return false, nil
	} else if err == nil {
		// file does exist
		file, err := os.Open(filename)
		defer file.Close()
		jsonParser := json.NewDecoder(file)
		gamelist := make([]Game, len(T.GameList), len(T.GameList))
		if err = jsonParser.Decode(&gamelist); err != nil {
			return false, err
		}
		// only load the completed games:
		for i, _ := range gamelist {
			if gamelist[i].Completed {
				T.GameList[i] = gamelist[i]
			}
		}
		return true, nil
	}
	return false, nil
}

func (T *Tourney) LoadFile(filename string) error {
	// Try to open the file:
	tourneyFile, err := os.Open(filename)
	defer tourneyFile.Close()
	if err != nil {
		fmt.Println("Error opening ", filename, ".", err.Error())
		return err
	}
	//Try to decode the file:
	jsonParser := json.NewDecoder(tourneyFile)
	if err = jsonParser.Decode(T); err != nil {
		fmt.Println("Error parsing .tourney file.", err.Error())
		return err
	}
	// Create the game list:
	T.GenerateMatchups()
	fmt.Println("Successfully Loaded ", filename)
	// Load the opening book:
	if T.BookLocation != "" {
		if err := LoadBook(T); err != nil {
			fmt.Println("Could not load opening book:", err)
		} else {
			fmt.Println("Successfully loaded opening book:", T.BookLocation)
			// TODO: this will have to be changed when more opening book support is added.
			if e := PlayOpenings(T); e != nil {
				fmt.Println("Error playing openings from book. ", e)
			} else {
				fmt.Println("Successfully applied openings from book.")
				/* DEBUG
				for i, _ := range T.GameList {
					fmt.Println(T.GameList[i].StartingFEN)
				}
				*/
			}
		}
	} else {
		fmt.Println("No opening book specified.")
	}
	// Check if this tourney was previously stopped midway
	if loaded, err := LoadPreviousResults(T); err != nil {
		return err
	} else if loaded {
		fmt.Println("Successfully loaded previous results.")
	}
	return nil
}

func (T *Tourney) LoadDefault() error {
	//TODO: I dont really like the name of this function
	var err error
	//Loads default.tourney
	if err = T.LoadFile("default.tourney"); err != nil {
		// something is wrong, so just load 40/2 CCLR settings:
		T.Event = "Tourney"
		T.Engines = make([]Engine, 0)
		T.TestSeats = 1
		T.Carousel = true
		T.Moves = 40
		T.Time = 120 //seconds
		T.Repeating = true
		T.Rounds = 30
		T.QuitAfter = false
	}
	return err
}

func (T *Tourney) GenerateMatchups() {
	// Deduce the needed information for the tourney to run.
	// This includes populating the game list.

	fmt.Println("Generating matchups between engines.")

	//Count the number of games:
	// TODO: VERIFY FORMULA!
	//S := T.TestSeats *( (T.TestSeats +1 )/2 ) // = Sum_{0}^{n} k
	//gameCount := T.Rounds * (T.TestSeats * len(T.Engines) - S)
	//T.GameList = make([]Game,gameCount)
	var def Game
	def.initialize()
	def.time = T.Time
	def.moves = T.Moves
	def.repeating = T.Repeating
	def.Completed = false

	// Non-Carousel:
	for t := 0; t < T.TestSeats; t++ {
		//Go around the test seats:
		for e := t + 1; e < len(T.Engines); e++ {
			//Now go around each opponent for that test seat:
			for r := 0; r < T.Rounds; r++ {
				//Finally all the rounds for that matchup:
				nextGame := def
				nextGame.Player[r%2] = T.Engines[t]
				nextGame.Player[(r+1)%2] = T.Engines[e]
				T.GameList = append(T.GameList, nextGame)
			}
		}
	}
	// TODO: Carousel
}

func GetResults(T *Tourney) []Record {
	var r []Record

	// helper function:
	indexOf := func(e Engine) int {
		for i, _ := range r {
			if r[i].Player.Name == e.Name {
				return i
			}
		}
		r = append(r, Record{Player: e})
		return len(r) - 1
	}
	// workhorse:
	for i, _ := range T.GameList {
		for color := WHITE; color <= BLACK; color++ {
			ind := indexOf(T.GameList[i].Player[color])
			if T.GameList[i].Completed {
				if winner := T.GameList[i].Result; winner == DRAW {
					r[ind].Draws++
				} else if winner == color {
					r[ind].Wins++
				} else {
					r[ind].Losses++
				}
			} else {
				r[ind].Incomplete++
			}
		}
	}
	return r
}

func ShowResults(T *Tourney) {
	// TODO: 	-add more detail. such as w-l-d for each matchup.

	// find the length of the longest name, for formatting purposes:
	var longestName int
	results := GetResults(T)
	for _, record := range results {
		if len(record.Player.Name) > longestName {
			longestName = len(record.Player.Name)
		}
	}
	for _, record := range results {
		fmt.Print(record.Player.Name, strings.Repeat(" ", longestName-len(record.Player.Name)), ":\t")
		fmt.Print(record.Wins, "-", record.Losses, "-", record.Draws, " remaining: ", record.Incomplete, "\t")
		score := float64(record.Wins) + 0.5*float64(record.Draws)
		possible := float64(record.Wins + record.Losses + record.Draws)
		fmt.Print(score, "/", possible, "\t")
		if possible > 0 {
			fmt.Printf("%.2f", 100*(score/possible))
			fmt.Print("%\n")
		}
	}
}
