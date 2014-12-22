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
 	-Vertical score graph. rows will be move#'s, cols will be the graph.
 	-ability to pipe commands

*******************************************************************************/

package main

import (
	"bufio"
	"fmt"
	"os"
	//"strconv"
	//"strings"
	"runtime"
	"sync"
)

type GlobalSettings struct {
	WorkerDirectory   string
	LogDirectory      string
	TemplateDirectory string
	SaveDirectory     string
	BookDirectory     string
	ServerPort        int
	WebPort           int
	EngineFilePort    int
}

func DefaultSettings() GlobalSettings {
	return GlobalSettings{
		WorkerDirectory:   "worker/",
		LogDirectory:      "logs/",
		TemplateDirectory: "templates/",
		SaveDirectory:     "",
		BookDirectory:     "",
		ServerPort:        9000,
		WebPort:           8080,
		EngineFilePort:    9001,
	}
}

var Settings GlobalSettings

func main() {

	// TODO: put settings in a file:
	Settings = DefaultSettings()

	fmt.Println("\nProject: Tourney Started\n")
	PrintSysStats()
	fmt.Println()
	// Check for a lanuch arguement with for a .tourney file
	// .tourney files contain all of the settings needed
	// to start a tourney without any terminal input.

	// validate that the file exists and is valid:

	// when no .tourney file is provided or is invalid, should load default.tourney
	/*
		var ActiveTourneys []*Tourney
		var SelectedIndex int

		def, _ := LoadDefault()
		ActiveTourneys = append(ActiveTourneys, def)

		SelectedIndex = 0
		ListActiveTourneys(ActiveTourneys, SelectedIndex)
	*/

	// TODO: Other launch arguements

	var Tourneys TourneyList
	def, _ := LoadDefault()
	Tourneys.Add(def)

	ListActiveTourneys(&Tourneys)

	// REPL:

	inputReader := bufio.NewReader(os.Stdin)

	var wg sync.WaitGroup
	var quit bool
	var prompt string
	for !quit {
		if len(Tourneys.List) > 0 {
			prompt = Tourneys.Selected().Event + "> "
		} else {
			prompt = "> "
		}
		fmt.Print(prompt)
		line, _ := inputReader.ReadString('\n')
		quit = Eval(line, &Tourneys, &wg)
	}
	wg.Wait()
	// DEBUG:
	fmt.Print("\nGoroutines: ", runtime.NumGoroutine(), "\n")
}

// Helper for common channel usage:
func blocks(c chan struct{}) bool {
	// until i can figure out a generic way to do this, this function will only support struct{} type
	select {
	case <-c:
		return false
	default:
	}
	return true
}
