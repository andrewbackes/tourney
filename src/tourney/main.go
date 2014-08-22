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
 	-Opening Book (non-pgn)
 	-Distributed game playing
 	-http output
 	-Vertical score graph. rows will be move#'s, cols will be the graph.
 	-ability to pipe commands

*******************************************************************************/

package main

import (
	"bufio"
	"fmt"
	"os"
	//"strconv"
<<<<<<< HEAD
	//"strings"
=======
>>>>>>> 7c94c876d91adedf33192d4b5aafc58cacb9af67
	"sync"
)

type Context struct {
	Done chan struct{}
}

func main() {

	fmt.Println("\nProject: Tourney Started\n")

	// Check for a lanuch arguement with for a .tourney file
	// .tourney files contain all of the settings needed
	// to start a tourney without any terminal input.

	// validate that the file exists and is valid:

	// when no .tourney file is provided or is invalid, should load default.tourney
	var ActiveTourneys []*Tourney
	var SelectedIndex int

	def, _ := LoadDefault()
	ActiveTourneys = append(ActiveTourneys, def)

	SelectedIndex = 0
	ListActiveTourneys(ActiveTourneys, SelectedIndex)
	// TODO: Other launch arguements

	// REPL:

	inputReader := bufio.NewReader(os.Stdin)

	var wg sync.WaitGroup
	var quit bool
	var prompt string
	for !quit {
		if len(ActiveTourneys) > 0 {
			prompt = ActiveTourneys[SelectedIndex].Event + "> "
		} else {
			prompt = "> "
		}
		fmt.Print(prompt)
		line, _ := inputReader.ReadString('\n')
		ActiveTourneys, quit = Eval(line, ActiveTourneys, &SelectedIndex, &wg)
	}
	fmt.Print("\n")
	wg.Wait()
}
