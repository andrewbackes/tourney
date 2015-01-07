/*

 Project: Tourney

 Module: commands
 Description: handles all of the commands from the user.

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/15/2014

TODO:
	- Need to prevent the user from using commands that both write to the same
	  object. Like RunTourney() and HostTourney()

	- delete command can have filepath issues.

*/

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

//
// Evaluate a command from the user.
// There is only one TourneyList and WaitGroup, it should be declared in main()
//
func Eval(command string, Tourneys *TourneyList, wg *sync.WaitGroup) bool {
	var queQuit bool
	command = strings.Trim(command, "\n")
	command = strings.Trim(command, " ")
	words := strings.Fields(command)
	T := Tourneys.Selected()

	// helper:
	type Command struct {
		label   []string
		desc    string
		format  string
		example string
		f       func()
	}

	PrintCommand := func(cmd Command) string {
		str := "Command:     "
		for _, c := range cmd.label {
			str += c + ", "
		}
		str = strings.Trim(str, ", ") + "\n"
		if cmd.format != "" {

			str += "How to use:  " + cmd.format + "\n"
		}
		if cmd.example != "" {
			str += "Example:     " + cmd.example + "\n"
		}
		str += "Description: " + cmd.desc
		return str
	}

	// Slice of possible commands, so we can search through them (or print them) later:
	var commands []Command
	commands = []Command{
		/*
			{
				label: []string{"hud"},
				desc:  "Displays the HUD for the current game. Prints the game board and other game state information.",
				f: func() {

				}},
		*/
		{
			label: []string{"info"},
			desc:  "Prints the configuration information for the selected tourney.",
			f: func() {
				Tourneys.Selected().Print()
			}},
		{
			label: []string{"settings"},
			desc:  "Displays Tourney's program settings.",
			f: func() {
				fmt.Print(Settings)
			}},
		{
			label: []string{"play", "start"},
			desc:  "Starts playing the currently selected tourney on the local machine. Use the 'stop' command to stop playing.",
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
						fmt.Println("To add an engine to play in this tournament type: addengine")
					}
					wg.Done()
				}()
				return
			}},
		{
			label: []string{"broadcast"},
			desc:  "Broadcasts the currently selected tourney over http port " + strconv.Itoa(Settings.WebPort) + ". Broadcasting is enabled by default. The default port is specified in the 'tourney.settings' file.",
			f: func() {
				if !Tourneys.broadcasting {
					fmt.Println("Broadcasting http on port " + strconv.Itoa(Settings.WebPort))
					go func() {
						//if err := Broadcast(&T, selected); err != nil {
						if err := Broadcast(Tourneys); err != nil {
							fmt.Println(err)
						} else {

						}
					}()
					Tourneys.broadcasting = true
				} else {
					fmt.Println("Already broadcasting on port " + strconv.Itoa(Settings.WebPort))
				}
				return
			}},
		{
			label: []string{"stop"},
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
			label: []string{"results"},
			desc:  "Displays the results of the currently selected tourney. Results may be displayed even if the tournament is incomplete.",
			f: func() {
				fmt.Print(SummarizeResults(Tourneys.Selected()))
				fmt.Println("To see more details type, 'games' or 'g'")
			}},
		{
			label: []string{"games"},
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
		{
			label: []string{"new"},
			desc:  "Creates a new tourney with default settings.",
			f: func() {
				fmt.Println("This command is not yet supported.")
				/*
					if N, err := LoadDefault(); err == nil {
						T = append(T, N)
						*selected = len(T) - 1
						ListActiveTourneys(T, *selected)
					}
				*/
			}},
		{
			label: []string{"ls"},
			desc:  "Displays a list of currently loaded tourneys.",
			f: func() {
				ListActiveTourneys(Tourneys)
			}},
		{
			label: []string{"quit", "q"},
			desc:  "Stops any running tournament after the next game and quits the program.",
			f: func() {
				if blocks(Tourneys.Selected().Done) {
					close(Tourneys.Selected().Done)
				}
				queQuit = true
			}},
		{
			label: []string{"help"},
			desc:  "Displays help for a specific command.",
			f: func() {
				var str string
				if len(words) > 1 {
					var cmd Command
					// find the command to help with:
					for _, c := range commands {
						for _, l := range c.label {
							if l == words[1] {
								cmd = c
							}
						}
					}
					if cmd.label == nil {
						fmt.Println("That is not a valid command. Use the 'commands' command to see a list of commands.")
						return
					}
					// display the help:
					str = PrintCommand(cmd)
				} else {
					str = "To get help with a specific command type: help <command>\nFor a list of commands type 'commands'"
				}
				fmt.Println(str)
			}},
		{
			label: []string{"commands"},
			desc:  "Displays a list of supported commands.",
			f: func() {
				cmds := ""
				for _, c := range commands {
					for i, _ := range c.label {
						cmds += c.label[i] + ", "
					}
				}
				cmds = strings.Trim(cmds, ", ")
				fmt.Println(cmds)
			}},
		{
			label: []string{"host"},
			desc:  "Hosts the currently loaded tournament on the local machine. However, does not start playing in the tournament.",
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
				} else {
					fmt.Println("No IP was specified. Assuming localhost.")
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
				// ex: buildbook "testing.pgn" 20 WhiteElo >2700 BlackElo >2700 Result =1/2-1/2
				var filename string
				// filename:
				if len(words) > 1 {
					filename = words[1]
				} else {
					fmt.Println("Please specify a filename.")
					return
				}
				// move count:
				moves := 12 //default
				if len(words) > 2 {
					moves, _ = strconv.Atoi(words[2])
				}
				// filters:
				var filters []PGNFilter
				for i := 3; i+1 < len(words); i = i + 2 {
					filters = append(filters, PGNFilter{Tag: words[i], Value: words[i+1]})
				}
				if b, err := LoadOrBuildBook(filename, moves, filters); err == nil {
					fmt.Println("Success.")
					fmt.Println(b.String())
				} else {
					fmt.Println("Failed:", err.Error())
				}
				return
			}},
		{
			label: []string{"rebuild", "rebuildbook"},
			desc:  "Rebuilds an opening book from a PGN.",
			f: func() {
				var filename string
				// filename:
				if len(words) > 1 {
					filename = words[1]
				}
				if err := os.Remove(filename[:len(filename)-3] + "book"); err != nil {
					fmt.Println(err)
				} else {
					Eval(command[2:], Tourneys, wg)
				}
				return
			}},
		{
			label: []string{"delete", "rm"},
			desc:  "Deletes all previous tourney data.",
			f: func() {
				extensions := []string{".data", ".pgn", ".txt"}
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
		{
			label: []string{"engine", "add", "addengine"},
			desc:  "Adds an engine to the tournament.",
			f: func() {
				inputReader := bufio.NewReader(os.Stdin)
				fmt.Println("Adding an engine to the tournament.")
				fmt.Print("Engine Name: ")
				name, _ := inputReader.ReadString('\n')
				fmt.Print("File Path: ")
				path, _ := inputReader.ReadString('\n')
				fmt.Print("Protocol (UCI or WINBOARD):")
				prot, _ := inputReader.ReadString('\n')
				Tourneys.Selected().AddEngine(strings.Trim(name, "\n"), strings.Trim(path, "\n"), strings.Trim(prot, "\n"))
				return
			}},
		{
			label:   []string{"timecontrol", "time", "settime"},
			format:  "timecontrol <moves>/<milliseconds>:<millisec added after each move>",
			example: "time 40/2000:10",
			desc:    "Modifies the time control of the tournament.",
			f: func() {
				if len(words) > 1 && strings.Contains(words[1], "/") {
					var m, t, i string
					m = strings.Split(words[1], "/")[0]
					t = strings.Split(words[1], "/")[1]
					if strings.Contains(t, ":") {
						i = strings.Split(t, ":")[1]
						t = strings.Split(t, ":")[0]
					}
					moves, _ := strconv.Atoi(m)
					time, _ := strconv.Atoi(t)
					inc, _ := strconv.Atoi(i)
					Tourneys.Selected().SetTimeControl(int64(moves), int64(time), int64(inc), Tourneys.Selected().Repeating)
				} else {
					fmt.Println("Not enough information to change time control.")
				}
			}},
		{
			label:   []string{"rounds"},
			format:  "rounds <#>",
			example: "rounds 10",
			desc:    "Modifies or Displays the number of rounds in the tournament.",
			f: func() {
				if len(words) > 1 {
					n, _ := strconv.Atoi(words[1])
					if n > 0 {
						T.SetRounds(n)
					} else {
						fmt.Println("Can not have 0 rounds.")
					}
				} else {
					fmt.Println("Rounds:", T.Rounds)
					fmt.Println("If you would like to change the number of rounds type: rounds <#>.")
				}
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
