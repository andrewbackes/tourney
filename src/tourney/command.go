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
	//TODO

	// This function is really really ugly!!!
	// Q: 	Can there be a map whos keys are the avaliable commands
	//		and the values are pointers to their individual functions?

	words := strings.Fields(command)
	if len(words) > 0 {
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
			err = showHelp()
		case "quit", "q":
			err = quit(T)
			quitFlag = true
		}
	}
	return
}

func commandLoop(tourney *Tourney) {
	// Root level command loop -
	// continuously accepts and executes the users commands.

	prompt := "Tourney> "
	inputReader := bufio.NewReader(os.Stdin)
	inputChan := make(chan string)

	quitChan := make(chan bool)
	defer close(quitChan)

	// UPSTREAM: ask user for a command
	fmt.Print(prompt)
	go func() {
		for {
			line, _ := inputReader.ReadString('\n')
			select {
			case inputChan <- line:
			case <-quitChan:
				break
			}
		}
		close(inputChan)
	}()

	// DOWNSTREAM: execute the command
	for i := range inputChan {
		go func() {
			quit, _ := doCommand(i, tourney)
			quitChan <- quit
		}()
		if <-quitChan {
			return
		}
		fmt.Print(prompt)
	}
	fmt.Print("\n")
}

func quit(T *Tourney) error {
	// Quit the program. So take care of business first.

	if T.State == RUNNING {
		T.Stop()
	}

	return nil
}

func showHelp() error {
	// TODO : What dan said!

	fmt.Println("Supported Commands:")
	fmt.Println("quit \t stops any running tournament and exits")
	fmt.Println("help \t displays this list")

	return nil
}
