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
 	-ability to pipe commands

*******************************************************************************/

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var Settings GlobalSettings

const SettingsFile = "tourney.settings"

// handleLaunchArgs processes the command line args
func handleLaunchArgs(repl, broadcast, loaddefault *bool, controller *Controller) {
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
			controller.Enque("load " + param)
			*loaddefault = false
		case "-play":
			// open and play the tourney:
			controller.Enque("load " + param)
			controller.Enque("start")
			*loaddefault = false
		case "-host":
			// open and host the tourney:
			controller.Enque("load " + param)
			controller.Enque("host")
			*loaddefault = false
		case "-connect":
			// connect to a host
			controller.Enque("connect " + param)
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

// changeWorkingDir adjusts working directory to where the program is located
func changeWorkingDir() {

	cd, er := filepath.Abs(filepath.Dir(os.Args[0]))
	err := os.Chdir(cd)
	if er != nil || err != nil {
		fmt.Println("Could not change working directory to ", cd)
	}
}

func main() {

	changeWorkingDir()

	title := "Tourney"
	fmt.Println("\n" + strings.Repeat(" ", (80-len(title))/2) + title + "\n")
	//PrintSysStats()

	loadSettings()

	controller := NewController()

	// Depending on the command line arguements, we will adjust these flags:
	repl := true
	broadcast := true
	loaddefault := true

	handleLaunchArgs(&repl, &broadcast, &loaddefault, &controller)

	if loaddefault {
		controller.Enque("loaddefault")
	}
	if repl {
		go ConsoleUI(&controller)
	}
	if broadcast {
		go WebUI(&controller)
		//controller.Enque("broadcast")
	}

	controller.Start()
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
