/*

 Project: Tourney

 Module: tourney
 Description: holds the tourney object and methods for interacting with it
 
 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*/

package main

import (
 	"fmt"
 	"os"
 	"encoding/json"
)

type State int

const (
	UNSTARTED State = iota
	RUNNING  				// in progress
	PAUSED					// incomplete but not running
	COMPLETED				// completed
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
	Moves uint8 //moves per time control
	Time uint64 // time per time control in seconds
	Repeating bool //restart time after moves hits

	Rounds uint8 //number of games each engine will play

	// Opening book information:

	QuitAfter bool //Quit after the tourney is complete.


// Control settings (Determined while tourney is running, or when the tourney starts)
	state State //flag to indicate: running, paused, stopped
	GameList []Game //list of all games in the tourney. should be populated when the tourney starts
	activeGame *Game //points to the currently running game in the list. Rethink this for multiple running games at a later time.

}

func (T *Tourney) LoadFile(filename string) error {
	
	// Try to open the file:
	tourneyFile, err := os.Open(filename)
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

func (T *Tourney) LoadDefault() {
//TODO: I dont really like the name of this function

	//Loads default.tourney
	err := T.LoadFile("default.tourney")

	if err != nil {
	//if the file doesnt exist or has invalid data, then load 40/2 CCLR settings:
		T.Name = "Tourney"
		T.Engines = make([]Engine,0)
		T.TestSeats = 1
		T.Carousel = true
		T.Moves = 40
		T.Time = 120 //seconds
		T.Repeating = true
		T.Rounds = 30
		T.QuitAfter = false
	}
}

