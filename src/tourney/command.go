/*

 Project: Tourney

 Module: commands
 Description: handles all of the commands from the user.

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/15/2014

*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// system wide commands should be: start, stop, pause, restart, new, quit, help

func doCommand(command string, T *Tourney) (quitFlag bool, err error) {
	// This function is really really ugly!!!
	// Q: 	Can there be a map whos keys are the avaliable commands
	//		and the values are pointers to their individual functions?

	words := strings.Fields(command)
	switch words[0] {
	case "start", "s":
		T.Start()
	case "stop", "p":
		T.Stop()
	case "new", "n":
		//TODO: prompt menu with options for new tourney
		err = T.LoadFile("default.tourney")
	case "load", "l":
		err = T.LoadFile(words[1])
	case "help", "h":
		quitFlag, err = showHelp()
	case "quit", "q":
		err = quit(T)
		quitFlag = true
	}
	return
}

func commandLoop(tourney *Tourney) {
	// Root level command loop -
	// continuously accepts and executes the users commands.
	// meant to be ran in the master thread
	var quitFlag bool
	var err error
	var line string

	prompt := "Tourney> "

	for !quitFlag {
		fmt.Print(prompt)
		//i, err2 := fmt.Scanln(&input)
		//fmt.Println("i: ", i, "err2: ", err2)

		input := bufio.NewReader(os.Stdin)
		line, err = input.ReadString('\n')

		quitFlag, err = doCommand(line, tourney)

		if err != nil {
			fmt.Println("An error occured with that command. You can type 'help' if you need.")
		}
	}
}

func quit(T *Tourney) error {
	// Quit the program. So take care of business first.

	if T.State == RUNNING {
		T.Stop()
	}

	return nil
}

func showHelp() (quit bool, err error) {
	fmt.Println("Supported Commands:")
	fmt.Println("quit \t stops any running tournament and exits")
	fmt.Println("help \t displays this list")

	quit = false
	return
}
