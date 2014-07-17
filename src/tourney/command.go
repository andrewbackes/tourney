/*

 Project: Tourney

 Module: commands
 Description: handles all of the commands from the user.
 
 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/15/2014

*/

package main

import (
 	"fmt"
 	//"strings"
)

// system wide commands should be: start, stop, pause, restart, new, quit, help

func doCommand(command string) (success bool, quitFlag bool){
	// This function is really really ugly!!!
	// Q: 	Can there be a map whos keys are the avaliable commands
	//		and the values are pointers to their individual functions?

	switch command {
		case "new", "n":

		case "help", "h":
			success, quitFlag = showHelp()
		case "quit", "q":
			success = true
			quitFlag = true
	}
	return
}

func commandLoop(tourney *Tourney) {
	// continuously accepts and executes the users commands.
	// meant to be ran in the master thread
	quitFlag := false
	successFlag := false
	prompt := "Tourney> "

	var input string
	for !quitFlag {
		fmt.Print(prompt)
		fmt.Scanf("%s", &input)
		
		successFlag, quitFlag = doCommand(input)
		
		if !successFlag {
			fmt.Println("Invalid Command. You can type 'help' if you need.")
		}
	}
}

func showHelp() (success, quit bool) {
	fmt.Println("Supported Commands:")
	fmt.Println("quit \t stops any running tournament and exits")
	fmt.Println("help \t displays this list")

	success = true
	quit = false
	return
}

