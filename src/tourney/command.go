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
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func Eval(command string, T []*Tourney, selected *int, wg *sync.WaitGroup) ([]*Tourney, bool) {
	var queQuit bool
	command = strings.Trim(command, "\n")
	command = strings.Trim(command, " ")
	words := strings.Fields(command)

	// helper:
	type Command struct {
		label []string
		desc  string
		f     func()
	}

	// Slice of possible commands, so we can search through them (or print them) later:
	var commands []Command
	commands = []Command{
		{
			label: []string{"display", "d"},
			desc:  "Displays the HUD for the current game.",
			f: func() {

			}},
		{
			label: []string{"tourney", "t"},
			desc:  "Prints the settings of the selected tourney.",
			f: func() {
				T[*selected].Print()
			}},
		{
			label: []string{"settings", "e"},
			desc:  "Changes the settings of the current tourney.",
			f: func() {
				Setup(T[*selected])
			}},
		{
			label: []string{"start", "s"},
			desc:  "Starts the currently selected tourney.",
			f: func() {
				go func() {
					T[*selected].Done = make(chan struct{})
				}()
				wg.Add(1)
				go func() {
					if err := RunTourney(T[*selected]); err != nil {
						fmt.Println(err)
					}
					wg.Done()
				}()
				return
			}},
		{
			label: []string{"stop", "p"},
			desc:  "Stops the tourney after the next game completes.",
			f: func() {
				wg.Add(1)
				go func() {
					close(T[*selected].Done)
					wg.Done()
				}()
			}},
		{
			label: []string{"results", "r"},
			desc:  "Displays the results of the currently selected tourney.",
			f: func() {
				fmt.Print(SummarizeResults(T[*selected]))
				fmt.Println("To see more details type, 'games' or 'g'")
			}},
		{
			label: []string{"games", "g"},
			desc:  "Displays the results of each game in the selected tourney.",
			f: func() {
				fmt.Print(SummarizeGames(T[*selected]))
			}},
		{
			label: []string{"load", "l"},
			desc:  "Loads a .tourney file.",
			f: func() {
				// loads the file and moves the selected tourney to the new one.
				filename := strings.Replace(command, words[0]+" ", "", 1)
				filename = strings.Trim(filename, "\r\n") // for windows
				filename = strings.Trim(filename, "\n")   // for *nix
				filename = strings.Replace(filename, ".tourney", "", 1) + ".tourney"
				if N, err := LoadFile(filename); err == nil {
					T = append(T, N)
					*selected = len(T) - 1
					ListActiveTourneys(T, *selected)
				}
			}},
		{
			label: []string{"new", "n"},
			desc:  "Creates a new tourney with default settings.",
			f: func() {
				if N, err := LoadDefault(); err == nil {
					T = append(T, N)
					*selected = len(T) - 1
					ListActiveTourneys(T, *selected)
				}
			}},
		{
			label: []string{"ls"},
			desc:  "Displays a list of currently loaded tourneys.",
			f: func() {
				ListActiveTourneys(T, *selected)
			}},
		{
			label: []string{"quit", "q"},
			desc:  "Quits the program",
			f: func() {
				queQuit = true
			}},
		{
			label: []string{"help", "h"},
			desc:  "Displays a menu of commands.",
			f: func() {
				fmt.Println("\nCommand:", "\t", "Description:\n")
				for _, c := range commands {
					for i, _ := range c.label {
						if i > 0 {
							fmt.Print(", ")
						}
						fmt.Print(c.label[i])
					}
					fmt.Println("\n\t", c.desc)
				}
				fmt.Println()
			}},
	}

	//helper:
	inSlice := func(l []string, s string) bool {
		for _, str := range l {
			if s == str {
				return true
			}
		}
		return false
	}
	// check if the user input is a valid command:
	if len(words) == 0 {
		return T, queQuit
	}
	for _, command := range commands {
		if inSlice(command.label, words[0]) {
			command.f()
			return T, queQuit
		}
	}
	// check to see if the user typed in an index or name of a tourney:
	var adjustedSelected bool
	for i, _ := range T {
		if (strconv.Itoa(i+1) == command) || (strings.ToLower(T[i].Event) == strings.ToLower(command)) {
			*selected = i
			adjustedSelected = true
			ListActiveTourneys(T, *selected)
		}
	}

	// No valid command was found:
	if !adjustedSelected {
		fmt.Print("There was an error processing that command. Type 'help' for a list of commands.\n")
	}
	return T, queQuit
}

func ListActiveTourneys(actT []*Tourney, selT int) {
	fmt.Print("\nThe following tourneys are currently loaded:\n")

	for i, _ := range actT {
		str := ""
		if i == selT {
			str += " --> "
		} else {
			str += "     "
		}
		str += strconv.Itoa(i+1) + ". " + actT[i].Event
		fmt.Println(str)
	}
	fmt.Println("\nTo select a different tourney from the list, type its name or number. ")
	fmt.Println("To load a tourney not listed, type 'load [filename]'")
	fmt.Println("To see this list again, type 'ls'")
	fmt.Println("For a list of additional commands, type 'help'\n")

}

func PrintSysStats() {
	fmt.Print("System: ", runtime.GOMAXPROCS(0), "/", runtime.NumCPU(), " CPUs.")
	fmt.Println()
}
