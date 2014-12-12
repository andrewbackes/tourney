/*******************************************************************************

 Project: Tourney

 Module: tourney

 Description: A Tourney object impliments the Game object. The Engine object is
 also implimented but is used for its data fields, its methods are ignored. The
 tournament takes place in the playLoop() method. The Start() and Stop() methods
 essentially modify the data feild "state" which is read by playLoop().

 TODO:
 	- Rename Tourney.Done to something more descriptive, like ForceQuit or
 	  something.
 	-Allow for games to be distributed to multiple machines to be played.
 		-Each machine will have to be benchmarked to determine equivalent
 		 Time control parameters.
 	-More tournament parameters
 	-Formatting results needs to be able to handle big numbers.
 	 Like: 35000-25000-10000
 	-Saving .tourney / .detail / .result / .pgn files when other already exist
	 should make a xxx1.xxx xxx2.xxx sort of thing.
	-Use text/template to save result files.
	-saved files should mimic the .tourney file name, not the event name.
	-previous tourney data file should be .data not .details

 BUGS:
 	-There may be an issue with things like: changing fields in the .tourney
 	 file when there is already a .details file. Because when the details are
 	 loaded, there may be a different number of games.
 	-Already played openings.
 	-Error handleing in RunTourney() incorrectly uses break

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
	"strconv"
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

	// Time control (Moves, Time, Repeating):
	Moves     int64 // Moves per Time control
	Time      int64 // Time per Time control in milliseconds
	BonusTime int64 // bonus Time added after each move
	Repeating bool  // restart Time after Moves hits

	Rounds int //number of games each engine will play

	// Opening book information:
	BookLocation string // File location of the book
	BookMoves    int    // Number of Moves to use out of the book
	BookPGN      []Game
	RandomBook   bool

	QuitAfter bool //Quit after the tourney is complete.

	// Control settings (Determined while tourney is running, or when the tourney starts)
	//State     Status //flag to indicate: running, paused, stopped
	//StateFlow chan Status
	//Flow     Context
	GameList []Game //list of all games in the tourney. populated when the tourney starts
	//activeGame *Game  //points to the currently running game in the list. Rethink this for multiple running games at a later Time.
	Done chan struct{}

	//For distribution:
	//GameQue          chan Game
	//CompletedGameQue chan Game
}

func RunTourney(T *Tourney) error {
	// TODO: verify that the settings currently loaded will not cause any problems.

	//var state Status
	for i, _ := range T.GameList {
		select {
		case <-T.Done:
			//channel closed, so stop.
			break
		default:
			//channel isnt closed, so keep playing
			fmt.Print("Round ", i+1, ": ", T.GameList[i].Player[WHITE].Name, " vs ", T.GameList[i].Player[BLACK].Name)
			if !T.GameList[i].Completed {
				fmt.Println("\nGame started.")
				fmt.Print("Playing from opening book... ")
				if err := PlayOpening(T, i); err != nil {
					fmt.Println("Failed:", err.Error())
					break
				}
				fmt.Println("Success.")
				if err := PlayGame(&T.GameList[i]); err != nil {
					fmt.Println(err.Error())
					break
				}
				fmt.Println("Game stopped.")
				T.GameList[i].PrintHUD()

				// Save progress:
				if err := Save(T); err != nil {
					return err
				}
			} else {
				fmt.Print(" -> ", []string{"1-0", "0-1", "1/2-1/2"}[T.GameList[i].Result], " - ", T.GameList[i].ResultDetail, "\n")
			}
		}
	}
	// Show results:
	fmt.Print(SummarizeResults(T))
	return nil
}

func Save(T *Tourney) error {
	// Save results:
	fmt.Print("Saving '" + T.Event + ".results'... ")
	if err := SaveResults(T); err != nil {
		fmt.Println("Failed.", err)
		//return err
	} else {
		fmt.Println("Success.")
	}
	// Save details:
	fmt.Print("Saving '" + T.Event + ".details'... ")
	if err := SaveDetails(T); err != nil {
		fmt.Println("Failed.", err)
		//return err
	} else {
		fmt.Println("Success.")
	}
	// Save PGN:
	fmt.Print("Saving '" + T.Event + ".pgn'... ")
	if err := SavePGN(T); err != nil {
		fmt.Println("Failed.", err)
		//return err
	} else {
		fmt.Println("Success.")
	}
	return nil
}

func SaveResults(T *Tourney) error {
	//check if the file exists:
	filename := T.Event + ".results"
	//var file *os.File
	//var err error
	if _, test := os.Stat(filename); os.IsNotExist(test) {
		// file doesnt exist
	} else if test == nil {
		// file does exist
		os.Remove(filename)
	}
	file, err := os.Create(filename)
	defer file.Close()
	summary := SummarizeResults(T) + SummarizeGames(T)
	_, err = file.WriteString(summary)

	return err
}

func SaveDetails(T *Tourney) error {
	//check if the file exists:
	filename := T.Event + ".details"
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
	//fmt.Println("Successfully saved " + filename)
	return nil
}

// Saves the completed games in pgn format:
func SavePGN(T *Tourney) error {
	//check if the file exists:
	filename := T.Event + ".pgn"
	if _, test := os.Stat(filename); os.IsNotExist(test) {
		// file doesnt exist
	} else if test == nil {
		// file does exist
		os.Remove(filename)
	}
	file, err := os.Create(filename)
	defer file.Close()

	var pgn string
	for i, _ := range T.GameList {
		if T.GameList[i].Completed {
			pgn += EncodePGN(&T.GameList[i])
		}
	}
	_, err = file.WriteString(pgn)
	return err
}

func LoadPreviousResults(T *Tourney) (bool, error) {
	filename := T.Event + ".details"
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

func LoadFile(filename string) (*Tourney, error) {

	// Try to open the file:
	fmt.Print("Loading tourney settings: '", filename, "'... ")
	tourneyFile, err := os.Open(filename)
	defer tourneyFile.Close()
	if err != nil {
		fmt.Println("Failed to open:", filename, ",", err.Error())
		return nil, err
	}
	// Make the object:
	T := new(Tourney)
	T.Done = make(chan struct{})

	// Try to decode the file:
	jsonParser := json.NewDecoder(tourneyFile)
	if err = jsonParser.Decode(T); err != nil {
		fmt.Println("Failed to decode:", err.Error())
		return nil, err
	}

	// Create the game list:
	T.GenerateGames()
	fmt.Print("Success.\n")

	// Load the opening book:
	if T.BookLocation != "" {
		fmt.Print("Loading opening book: '", T.BookLocation, "'... ")
		if err := LoadBook(T); err != nil {
			fmt.Println("Failed to load opening book:", err)
			return nil, err
		} else {
			fmt.Println("Success (", len(T.BookPGN), "Openings ).")
		}
	} else {
		fmt.Println("No opening book specified.")
	}
	// Check if this tourney was previously stopped midway
	fmt.Print("Loading previous tourney data... ")
	if loaded, err := LoadPreviousResults(T); err != nil {
		fmt.Println("Failed.", err)
		return nil, err
	} else if loaded {
		fmt.Println("Success.")
	} else {
		fmt.Println("Nothing to load.")
	}

	// verify engines:
	// BUG: gamelist is already generated at this point.
	/*
		fmt.Print("Verifying file integrity of engines... ")
		if err := T.VerifyEngineIntegrity(); err != nil {
			fmt.Println("Failed:", err)
		} else {
			fmt.Println("Success.")
		}
	*/

	return T, nil
}

func LoadDefault() (*Tourney, error) {
	//TODO: I dont really like the name of this function
	var err error
	// Create the object:
	var T *Tourney
	//Loads default.tourney
	if T, err = LoadFile("default.tourney"); err != nil {

		// something is wrong, so just load 40/2 CCLR settings:

		T = new(Tourney)
		T.Done = make(chan struct{})
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
	return T, err
}

func (T *Tourney) GenerateGames() {
	// Populates the game list with a generic unstarted game based
	// on the Tourney parameters.

	//fmt.Println("Generating matchups between engines.")

	//Count the number of games:
	// TODO: VERIFY FORMULA!
	//S := T.TestSeats *( (T.TestSeats +1 )/2 ) // = Sum_{0}^{n} k
	//gameCount := T.Rounds * (T.TestSeats * len(T.Engines) - S)
	//T.GameList = make([]Game,gameCount)
	var def Game
	def.initialize()
	def.Event = T.Event
	def.Time = T.Time
	def.Moves = T.Moves
	def.Repeating = T.Repeating
	def.Completed = false
	def.resetTimeControl()
	def.Board.Reset()
	def.CastleRights = [2][2]bool{{true, true}, {true, true}}
	def.EnPassant = 64
	def.Completed = false

	for t := 0; t < T.TestSeats; t++ {
		//Go around the test seats:
		if T.Carousel {
			for r := 0; r < T.Rounds; r = r + []int{2, 1}[T.Rounds%2] {
				for e := t + 1; e < len(T.Engines); e++ {
					nextGame := def
					nextGame.Player[r%2] = T.Engines[t]
					nextGame.Player[(r+1)%2] = T.Engines[e]
					T.GameList = append(T.GameList, nextGame)
					if T.Rounds%2 == 0 {
						nextNextGame := def
						nextNextGame.Player[r%2] = T.Engines[e]
						nextNextGame.Player[(r+1)%2] = T.Engines[t]
						T.GameList = append(T.GameList, nextNextGame)
					}
				}
			}
		} else {
			// Non-Carousel:
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
	}
	// Set the round numbers:
	for i, _ := range T.GameList {
		T.GameList[i].Round = i + 1
	}
}

// Print the settings of the tourney:
func (T *Tourney) Print() {
	// TODO: test seats
	var summary string
	summary = strings.Repeat("=", 80) + "\n Tourney Settings:\n" + strings.Repeat("=", 80) + "\n"
	summary += " Event:          " + T.Event + "\n"
	summary += " Site:           " + T.Site + "\n"
	summary += " Date:           " + T.Date + "\n"
	summary += " Rounds:         " + strconv.Itoa(T.Rounds) + "\n"
	summary += " Gauntlet Seats: " + strconv.Itoa(T.TestSeats) + "\n"
	// TODO: add remaining Time for non Repeating. use ':'
	summary += " Time control:   " + strconv.FormatInt(T.Moves, 10) + "/" +
		strconv.FormatInt(T.Time, 10) + " +" +
		strconv.FormatInt(T.BonusTime, 10) + "\n\n"

	summary += strings.Repeat("-", 80) + "\n Participants:\n" + strings.Repeat("-", 80) + "\n"
	for _, e := range T.Engines {
		summary += " Name:      " + e.Name + "\n  Path:     " + e.Path + "\n  Protocol: " + e.Protocol + "\n"
	}
	summary += "\n"
	summary += strings.Repeat("-", 80) + "\n Opening Book:\n" + strings.Repeat("-", 80) + "\n"
	summary += " Path:         " + T.BookLocation + "\n"
	summary += " # Book Moves: " + strconv.Itoa(T.BookMoves) + "\n"
	summary += " Randomize:    " + strconv.FormatBool(T.RandomBook) + "\n"

	fmt.Println(summary)
}

// Figures out if the tourney is complete:
func (T *Tourney) Complete() bool {
	for i, _ := range T.GameList {
		if !T.GameList[i].Completed {
			return false
		}
	}
	return true
}

func (T *Tourney) VerifyEngineIntegrity() error {
	// After the tourney is loaded and the engine paths are defined,
	// this goes through and MD5 checks the engine files against what was
	// previously ran. If it is a fresh tourney, then it saves the MD5
	// sums of the engines to check against later.

	for i, _ := range T.Engines {
		if err := T.Engines[i].ValidateEngineFile(); err != nil {
			return errors.New(fmt.Sprint(T.Engines[i].Path, "-", err))
		}
	}

	return nil
}
