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
	"sync"
)

// system wide commands should be: start, stop, pause, restart, new, quit, help

func doCommand(command string, T *Tourney, wg *sync.WaitGroup) bool {
	//TODO

	// This function is really really ugly!!!
	// Q: 	Can there be a map whos keys are the avaliable commands
	//		and the values are pointers to their individual functions?

	words := strings.Fields(command)
	if len(words) > 0 {
		switch words[0] {
		case "start", "s":
			T.StateFlow = make(chan Status)
			go func() {
				T.StateFlow <- RUNNING // BUG: Doing this more than once, this will block.
			}()
			wg.Add(1)
			go func() {
				if err := RunTourney(T); err != nil {
					fmt.Println(err)
				}
				wg.Done()
			}()
		case "stop", "p":
			wg.Add(1)
			go func() {
				T.StateFlow <- STOPPED // BUG: Doing this more than once, this will block.
				wg.Done()
			}()
		case "new", "n":
			//TODO: prompt menu with options for new tourney
			T.LoadFile("default.tourney") //TODO: this will cause unknown errors with the goroutine running the tourney
		case "load", "l":
			T.LoadFile(words[1]) //TODO: this will cause unknown errors with the goroutine running the tourney
		case "help", "h":
			showHelp()
		case "quit", "q":
			wg.Add(1)
			go func() {
				//close(T.StateFlow)
				//T.StateFlow <- STOPPED
				wg.Done()
			}()
			return true
		default:
			fmt.Print("There was an error processing that command. Type 'help' for a list of commands.\n")
		}
	}
	return false
}

func Controller(tourney *Tourney) {
	// continuously accepts and executes the users commands.

	prompt := "Tourney> "
	inputReader := bufio.NewReader(os.Stdin)

	var wg sync.WaitGroup
	var quit bool

	for !quit {
		fmt.Print(prompt)
		line, _ := inputReader.ReadString('\n')
		quit = doCommand(line, tourney, &wg)
	}
	fmt.Print("\n")
	wg.Wait()
}

func showHelp() error {
	// TODO : What dan said!

	fmt.Println("Supported Commands:")
	fmt.Println("quit \t stops any running tournament and exits")
	fmt.Println("help \t displays this list")

	return nil
}
