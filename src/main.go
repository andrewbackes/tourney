/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes
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

 BUG:
 	-current directory could differ from the directory the executable is in.

*******************************************************************************/

package main

import (
	"bufio"
	"fmt"
	"os"
	//"strconv"
	"strings"
	//"runtime"
	"path/filepath"
	"sync"
)

var Settings GlobalSettings

const SettingsFile = "tourney.settings"

// handleLaunchArgs processes the command line args 
func handleLaunchArgs(repl, broadcast, loaddefault *bool, Tourneys *TourneyList, wg *sync.WaitGroup) {
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
			Eval("load "+param, Tourneys, wg)
			*loaddefault = false
		case "-play":
			// open and play the tourney:
			Eval("load "+param, Tourneys, wg)
			Eval("start", Tourneys, wg)
			*loaddefault = false
		case "-host":
			// open and host the tourney:
			Eval("load "+param, Tourneys, wg)
			Eval("host", Tourneys, wg)
			*loaddefault = false
		case "-connect":
			// connect to a host
			Eval("connect "+param, Tourneys, wg)
			*repl = false
			*broadcast = false
			*loaddefault = false
		}
	}
}

// loadSettings loads Global Settings:
func loadSettings() {
	if err := Settings.Load(SettingsFile); err != nil {
		fmt.Println(err)
		fmt.Println("Using default program settings.")
		Settings = DefaultSettings()
		Settings.Save("tourney.settings")
	}
}

func main() {
	// Adjust working directory:
	cd, er := filepath.Abs(filepath.Dir(os.Args[0]))
	err := os.Chdir(cd)
	if er != nil || err != nil {
		fmt.Println("Could not change working directory to ", cd)
	}

	title := "Tourney"
	fmt.Println("\n" + strings.Repeat(" ", (80-len(title))/2) + title + "\n")
	//PrintSysStats()

	loadSettings()
	

	var wg sync.WaitGroup
	var Tourneys TourneyList
	// Depending on the command line arguements, we will adjust these flags:
	repl := true
	broadcast := true
	loaddefault := true

	handleLaunchArgs(&repl, &broadcast, &loaddefault, &Tourneys, &wg)

	if loaddefault {
		def, _ := LoadDefault()
		Tourneys.Add(def)
		//ListActiveTourneys(&Tourneys)
		Eval("ls", &Tourneys, &wg)
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
