/*

 Project: Tourney

 Module: tourney
 Description: holds the tourney object and methods for interacting with it

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*/

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

func (T *Tourney) playLoop() error {
	for T.State == RUNNING {
		for _, g := range T.GameList {
			g.Start() // TODO: adjust for partially completed tourneys/games
		}

		//Temporary:
		T.Stop()
	}
	return nil
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
