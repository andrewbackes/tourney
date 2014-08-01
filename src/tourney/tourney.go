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
	"encoding/json"
	"fmt"
	"os"
)

type Status int

const (
	UNSTARTED Status = iota
	RUNNING          // in progress
	STOPPED
)

type Tourney struct {
	// Predetermined Settings for a tourney:
	Name string //identifier for this tournament. Unique may be better?

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

	QuitAfter bool //Quit after the tourney is complete.

	// Control settings (Determined while tourney is running, or when the tourney starts)
	State      Status //flag to indicate: running, paused, stopped
	GameList   []Game //list of all games in the tourney. should be populated when the tourney starts
	activeGame *Game  //points to the currently running game in the list. Rethink this for multiple running games at a later time.

}

func (T *Tourney) playLoop() error {
	for T.State == RUNNING {
		for i := 0; i < len(T.GameList); i++ { // doing for i, g := range T.GameList makes a copy of the game. is there a way to point instead?
			T.GameList[i].Start() // TODO: adjust for partially completed tourneys/games
		}

		//Temporary:
		T.Stop()
	}
	// Show results:
	T.Status()

	return nil
}

func (T *Tourney) LoadFile(filename string) error {

	// Try to open the file:

	tourneyFile, err := os.Open(filename)
	// when i try to combine this with the if statement, i get an error about tourneyFile not being defined
	if err != nil {
		printError("opening .tourney file", err.Error())
		return err
	}

	//Try to decode the file:
	jsonParser := json.NewDecoder(tourneyFile)
	if err = jsonParser.Decode(T); err != nil {
		printError("parsing .tourney file", err.Error())
		return err
	}

	fmt.Println("Successfully Loaded ", filename)
	return nil
}

func (T *Tourney) LoadDefault() error {
	//TODO: I dont really like the name of this function
	var err error
	//Loads default.tourney
	if err = T.LoadFile("default.tourney"); err != nil {
		// something is wrong, so just load 40/2 CCLR settings:
		T.Name = "Tourney"
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

func (T *Tourney) Generate() {
	// Deduce the needed information for the tourney to run.
	// This includes populating the game list.

	fmt.Println("Generating Tourney Parameters.")

	//Count the number of games:
	// TODO: VERIFY FORMULA!
	//S := T.TestSeats *( (T.TestSeats +1 )/2 ) // = Sum_{0}^{n} k
	//gameCount := T.Rounds * (T.TestSeats * len(T.Engines) - S)
	//T.GameList = make([]Game,gameCount)
	def := Game{time: T.Time, moves: T.Moves, repeating: T.Repeating, state: UNSTARTED, completed: false}

	// Non-Carousel:
	for t := 0; t < T.TestSeats; t++ {
		//Go around the test seats:
		for e := t + 1; e < len(T.Engines); e++ {
			//Now go around each opponent for that test seat:
			for r := 0; r < T.Rounds; r++ {
				//Finally all the rounds for that matchup:
				nextGame := def
				nextGame.player[r%2] = T.Engines[t]
				nextGame.player[(r+1)%2] = T.Engines[e]
				T.GameList = append(T.GameList, nextGame)
			}
		}
	}
}

func (T *Tourney) Start() error {
	// Controls the state of the tourney.
	if T.State == UNSTARTED {
		T.Generate()
	}
	// TODO: error check to make sure it is safe to start it up right now

	T.State = RUNNING
	fmt.Println("Tourney is running.")

	// Begin playing games:
	T.playLoop()

	return nil
}

func (T *Tourney) Stop() error {
	// Controls the state of the tourney.
	// TODO: change the state of each game to STOPPED

	if T.State == RUNNING {
		for _, g := range T.GameList {
			if g.state == RUNNING {
				g.Stop()
			}
		}
		T.State = STOPPED
		fmt.Println("Tourney stopped.")
	}
	return nil
}

func (T *Tourney) Status() {
	// TODO: 	-add more detail. such as w-l-d for each matchup.
	//			-formatting is messed up for long and then short engine names
	type record struct {
		wins, losses, draws, remaining int
	}

	records := make(map[string]*record)
	for _, g := range T.GameList {
		// golang trick so that i can do map[key].field :
		if records[g.player[WHITE].Name] == nil {
			records[g.player[WHITE].Name] = &record{0, 0, 0, 0}
		}
		if records[g.player[BLACK].Name] == nil {
			records[g.player[BLACK].Name] = &record{0, 0, 0, 0}
		}

		if g.completed == false {
			records[g.player[WHITE].Name].remaining += 1
			records[g.player[BLACK].Name].remaining += 1
			continue
		}
		if g.result == DRAW {
			records[g.player[WHITE].Name].draws += 1
			records[g.player[BLACK].Name].draws += 1
		} else {
			records[g.player[g.result].Name].wins += 1
			records[g.player[[]Color{BLACK, WHITE}[g.result]].Name].losses += 1
		}
	}
	for _, e := range T.Engines {

		fmt.Print(e.Name, ": \t", records[e.Name].wins, "-", records[e.Name].losses, "-", records[e.Name].draws, ".\t")
		fmt.Print(float64(records[e.Name].wins) + (0.5 * float64(records[e.Name].draws)))
		gamesPlayed := records[e.Name].wins + records[e.Name].losses + records[e.Name].draws
		fmt.Print("/", gamesPlayed)

		if gamesPlayed > 0 {
			fmt.Print("\t")
			fmt.Printf("   %.2f", 100*(float64(records[e.Name].wins)+(0.5*float64(records[e.Name].draws)))/float64(gamesPlayed))
			fmt.Print("%")
		} else {
			fmt.Print("\t--.--%")
		}

		fmt.Print("\tRemaining: ", records[e.Name].remaining, "\n")

	}
}
