/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/14/2014

 Description: Plays tournaments between chess engines.

 The commands module recieves commands from the user.
 Once a tournament is started:

 A Tourney object impliments multiple Game objects. Each Game object impliments
 two Engine objects. The Engine object communicates with the chess engines
 through stdio. The Game object plays the games of the tournament with the
 playLoop() method.

 TODO:
 	-Opening Book (make sure to note the first moves out of the book and FEN)
 	-Distributed game playing
 	-http output
 	-Vertical score graph. rows will be move#'s, cols will be the graph.

*******************************************************************************/

package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	/*
		var G Game
		G.MakeMove(Move{Algebraic: "e2e4"})
		G.MakeMove(Move{Algebraic: "d7d5"})
		G.MakeMove(Move{Algebraic: "b1c3"})
		G.MakeMove(Move{Algebraic: "b8c6"})
		G.Completed = true
		G.Result = WHITE
		G.Event = "E"
		fmt.Print(EncodePGN(&G))
		return
	*/
	/*
		buf, err := ioutil.ReadFile("sample.pgn")
		if err != nil {
			fmt.Println(err)
		}
		s := string(buf)
		g := DecodePGN(s)
		for _, s := range g {
			fmt.Println(s, "\n\n\n\n")

		}
		return
	*/
	fmt.Println("Project: Tourney Started")

	// Until there is a need to have multiple Tourney objects to run at once,
	// this single object will just be passed around and manipulated:
	var tourney Tourney
	tourney.StateFlow = make(chan Status)
	// Check for a lanuch arguement with for a .tourney file
	// .tourney files contain all of the settings needed
	// to start a tourney without any terminal input.

	// validate that the file exists and is valid:

	// when no .tourney file is provided or is invalid, should load default.tourney
	tourney.LoadDefault()

	// TODO: Other launch arguements

	// and either go to the menu or the command loop
	Controller(&tourney)

}
