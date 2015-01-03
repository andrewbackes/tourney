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
	"strings"
	//"runtime"
	"sync"
)

type GlobalSettings struct {
	WorkerDirectory   string
	LogDirectory      string
	TemplateDirectory string
	SaveDirectory     string
	BookDirectory     string

	ServerPort     int
	WebPort        int
	EngineFilePort int

	MaxConnectionAttempts int
}

func DefaultSettings() GlobalSettings {
	return GlobalSettings{
		WorkerDirectory:       "worker/",
		LogDirectory:          "logs/",
		TemplateDirectory:     "templates/",
		SaveDirectory:         "",
		BookDirectory:         "",
		ServerPort:            9000,
		WebPort:               8080,
		EngineFilePort:        9001,
		MaxConnectionAttempts: 3,
	}
}

var Settings GlobalSettings

func main() {

	// TODO: put settings in a file:
	Settings = DefaultSettings()

	title := "Tourney"
	fmt.Println("\n" + strings.Repeat(" ", (80-len(title))/2) + title + "\n")
	//PrintSysStats()
	//fmt.Println()

	var wg sync.WaitGroup
	var Tourneys TourneyList
	// Depending on the command line arguements, we will adjust these flags:
	repl := true
	broadcast := true
	loaddefault := true

	// Parse the command line arguements:
	args := os.Args[1:]
	for i, arg := range args {
		if len(arg) == 2 {
			arg = strings.Replace(arg, "-o", "-open", -1)
			arg = strings.Replace(arg, "-c", "-connect", -1)
			arg = strings.Replace(arg, "-h", "-host", -1)
			arg = strings.Replace(arg, "-p", "-play", -1)
		}
		param := ""
		if len(args) > i+1 {
			param = args[i+1]
		}
		switch arg {
		case "-open":
			// open .tourney file
			Eval("load "+param, &Tourneys, &wg)
			loaddefault = false
		case "-play":
			// open and play the tourney:
			Eval("load "+param, &Tourneys, &wg)
			Eval("start", &Tourneys, &wg)
			loaddefault = false
		case "-host":
			// open and host the tourney:
			Eval("load "+param, &Tourneys, &wg)
			Eval("host", &Tourneys, &wg)
			loaddefault = false
		case "-connect":
			// connect to a host
			Eval("connect "+param, &Tourneys, &wg)
			repl = false
			broadcast = false
			loaddefault = false
		}
	}

	if loaddefault {
		def, _ := LoadDefault()
		Tourneys.Add(def)
		ListActiveTourneys(&Tourneys)
	}

	// Start web services.
	if broadcast {
		Eval("broadcast", &Tourneys, &wg)
	}

	// REPL:
	if repl {

		inputReader := bufio.NewReader(os.Stdin)

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
	}
	wg.Wait()
}

// Helper for common channel usage:
func blocks(c chan struct{}) bool {
	// until i can figure out a painless generic way to do this, this function will only support struct{} type
	select {
	case <-c:
		return false
	default:
	}
	return true
}
