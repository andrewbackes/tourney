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
)

func main() {
	/*
		var G Game
		G.LoadFEN("r3k2r/8/8/8/8/8/8/4K3 w kq - 0 1")
		G.MakeMove(Move{algebraic: "e1f1"})
		G.MakeMove(Move{algebraic: "a8a1"})
		G.PrintHUD()
	*/
	/*
		fmt.Println(G.FEN())

		G.MakeMove(Move{algebraic: "e1f1"})
		G.MakeMove(Move{algebraic: "h8h3"})
		G.MakeMove(Move{algebraic: "f1g1"})
		G.Print()
		fmt.Println(G.FEN())

		//divide(G, 1)
	*/
	/*
		for i := PAWN; i <= KING; i++ {
			bitprint(G.board.pieceBB[WHITE][i])
			fmt.Println(" ")
			bitprint(G.board.pieceBB[BLACK][i])
			fmt.Println("----------------------------------------")
		}
	*/
	//G.Print()
	//PerftSuite("/Users/Andrew/Projects/tourney/bin/perftsuite.epd", 7) //passed up to FEN 30.

	/*
		goal := []int64{1, 20, 400, 8902, 197281, 4865609, 119060324, 3195901860}

		for i := 1; i <= 6; i++ {
			nodes, chk, cstl, m, cap, prom, enpas := perft(G, i)
			fmt.Println("PERFT ", i, ": ", nodes, cap, enpas, cstl, prom, chk, m)

			if diff := int64(nodes) - int64(goal[i]); diff != 0 {
				fmt.Println("\toff by ", diff)
			}

		}
	*/
	//return
	fmt.Println("Project: Tourney Started")

	// Until there is a need to have multiple Tourney objects to run at once,
	// this single object will just be passed around and manipulated:
	var tourney Tourney

	// Check for a lanuch arguement with for a .tourney file
	// .tourney files contain all of the settings needed
	// to start a tourney without any terminal input.

	// validate that the file exists and is valid:

	// when no .tourney file is provided or is invalid, should load default.tourney
	tourney.LoadDefault()

	// TODO: Other launch arguements

	// and either go to the menu or the command loop
	commandLoop(&tourney)

}

type OutputType int

const (
	STANDARD OutputType = iota
	ERROR
	ENGINE_INPUT
	ENGINE_OUTPUT
)

func Output(a ...interface{}) {

}
