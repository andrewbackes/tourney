/*

 Project: Tourney

 Module: commands
 Description: handles all of the commands from the user.

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/15/2014

TODO:
	-Need to prevent the user from using commands that both write to the same
	 object. Like RunTourney() and HostTourney()

*/

package main

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

//func Eval(command string, T []*Tourney, selected *int, wg *sync.WaitGroup) ([]*Tourney, bool) {
func Eval(command string, Tourneys *TourneyList, wg *sync.WaitGroup) bool {
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
				//T[*selected].Print()
				Tourneys.Selected().Print()
			}},
		{
			label: []string{"settings", "e"},
			desc:  "Changes the settings of the current tourney.",
			f: func() {
				//Setup(T[*selected])
			}},
		{
			label: []string{"start", "s"},
			desc:  "Starts the currently selected tourney.",
			f: func() {
				go func() {
					//T[*selected].Done = make(chan struct{})
					Tourneys.Selected().Done = make(chan struct{})
				}()
				wg.Add(1)
				go func() {
					//if err := RunTourney(T[*selected]); err != nil {
					if err := RunTourney(Tourneys.Selected()); err != nil {
						fmt.Println(err)
					}
					wg.Done()
				}()
				return
			}},
		{
			label: []string{"broadcast", "b"},
			desc:  "Broadcasts the currently selected tourney over http port 8000.",
			f: func() {
				fmt.Println("Broadcasting http on port 8080.")
				go func() {
					//if err := Broadcast(&T, selected); err != nil {
					if err := Broadcast(Tourneys); err != nil {
						fmt.Println(err)
					}
				}()
				return
			}},
		{
			label: []string{"stop", "p"},
			desc:  "Stops the tourney after the next game completes.",
			f: func() {
				wg.Add(1)
				go func() {
					if blocks(Tourneys.Selected().Done) {
						close(Tourneys.Selected().Done)
					}
					wg.Done()
				}()
			}},
		{
			label: []string{"results", "r"},
			desc:  "Displays the results of the currently selected tourney.",
			f: func() {
				fmt.Print(SummarizeResults(Tourneys.Selected()))
				fmt.Println("To see more details type, 'games' or 'g'")
			}},
		{
			label: []string{"games", "g"},
			desc:  "Displays the results of each game in the selected tourney.",
			f: func() {
				fmt.Print(SummarizeGames(Tourneys.Selected()))
			}},
		{
			label: []string{"load", "l"},
			desc:  "Loads a .tourney file.",
			f: func() {
				// loads the file and Moves the selected tourney to the new one.
				filename := strings.Replace(command, words[0]+" ", "", 1)
				filename = strings.Trim(filename, "\r\n") // for windows
				filename = strings.Trim(filename, "\n")   // for *nix
				filename = strings.Replace(filename, ".tourney", "", 1) + ".tourney"
				if N, err := LoadFile(filename); err == nil {
					//T = append(T, N)
					//*selected = len(T) - 1
					Tourneys.Add(N)
					ListActiveTourneys(Tourneys)
				}
			}},
		/*
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
		*/
		{
			label: []string{"ls"},
			desc:  "Displays a list of currently loaded tourneys.",
			f: func() {
				ListActiveTourneys(Tourneys)
			}},
		{
			label: []string{"quit", "q"},
			desc:  "Quits the program",
			f: func() {
				if blocks(Tourneys.Selected().Done) {
					close(Tourneys.Selected().Done)
				}
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
		{
			label: []string{"host", "o"},
			desc:  "Hosts and runs a tourney.",
			f: func() {
				wg.Add(1)
				go func() {
					Tourneys.Selected().Done = make(chan struct{})
					if err := HostTourney(Tourneys.Selected()); err != nil {
						fmt.Println(err)
					}
					wg.Done()
				}()
				return
			}},
		{
			label: []string{"connect", "c"},
			desc:  "Connects to a host running a tourney.",
			f: func() {
				address := "127.0.0.1"
				if len(words) > 1 {
					address = words[1]
				}
				if !strings.Contains(address, ":") {
					address += fmt.Sprint(":", Settings.ServerPort)
				}
				wg.Add(1)
				go func() {
					ConnectAndWait(address)
					wg.Done()
				}()
				return
			}},
		{
			label: []string{"dowork"},
			desc:  "Connects to dirty-bit.com and becomes a worker.",
			f: func() {
				wg.Add(1)
				go func() {
					WorkForDirtyBit()
					wg.Done()
				}()
				return
			}},
		{
			label: []string{"build", "buildbook"},
			desc:  "Build an opening book from a PGN.",
			f: func() {
				var filename string
				if len(words) > 1 {
					filename = words[1]
				} else {
					fmt.Println("Please specify a filename.")
					return
				}
				fmt.Print("Building opening book from '", filename, "'...\n")
				if b, err := BuildBook(filename, 4); err == nil {
					fmt.Println("Success.")
					fmt.Println(*b)
				} else {
					fmt.Println("Failed:", err.Error())
				}
				return
			}},
		{
			label: []string{"delete", "rm"},
			desc:  "Deletes all previous tourney data.",
			f: func() {
				extensions := []string{".data", ".pgn", ".results"}
				for _, v := range extensions {
					if err := os.Remove(Tourneys.Selected().filename + v); err != nil {
						fmt.Println(err)
					} else {
						fmt.Println("Deleted " + Tourneys.Selected().filename + v)
					}
				}
				for i, _ := range Tourneys.Selected().GameList {
					path := "logs/" + Tourneys.Selected().Event + " round " + strconv.Itoa(i+1) + ".log"
					if err := os.Remove(path); err != nil {
						fmt.Println(err)
						break
					} else {
						fmt.Println("Deleted " + path)
					}
				}
				return
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
		return queQuit
	}
	for _, command := range commands {
		if inSlice(command.label, words[0]) {
			command.f()
			return queQuit
		}
	}
	// check to see if the user typed in an index or name of a tourney:
	var adjustedSelected bool
	for i, _ := range Tourneys.List {
		if (strconv.Itoa(i+1) == command) || (strings.ToLower(Tourneys.List[i].Event) == strings.ToLower(command)) {
			//*selected = i
			Tourneys.Index = i
			adjustedSelected = true
			ListActiveTourneys(Tourneys)
		}
	}

	// No valid command was found:
	if !adjustedSelected {
		fmt.Print("There was an error processing that command. Type 'help' for a list of commands.\n")
	}
	return queQuit
}

//func ListActiveTourneys(actT []*Tourney, selT int) {
func ListActiveTourneys(Tourneys *TourneyList) {
	fmt.Print("\nThe following tourneys are currently loaded:\n")

	for i, _ := range Tourneys.List {
		str := ""
		if i == Tourneys.Index {
			str += " --> "
		} else {
			str += "     "
		}
		str += strconv.Itoa(i+1) + ". " + Tourneys.List[i].Event
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
