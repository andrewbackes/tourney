/*******************************************************************************

 Project: 		Tourney
 Module: 		command
 Created: 		7/15/2014
 Author(s): 	Andrew Backes
 Description: 	Handles all of the possible user commands

TODO:
	- Need to prevent the user from using commands that both write to the same
	  object. Like RunTourney() and HostTourney()
	- delete command can have filepath issues.

*******************************************************************************/

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

type UserCommand struct {
	label          []string
	desc           string
	usage          string
	example        string
	f              func()
	classification int
}

const UI = 0
const TOURNEY_CONTROL = 1

func (C UserCommand) String() string {
	str := "Command(s):  "
	for _, c := range C.label {
		str += c + ", "
	}
	str = strings.Trim(str, ", ") + "\n"
	if C.usage != "" {

		str += "Usage:       " + C.usage + "\n"
	}
	if C.example != "" {
		str += "Example:     " + C.example + "\n"
	}
	str += "Description: " + C.desc
	return str
}

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

	// Slice of possible commands, so we can search through them (or print them) later:
	var commands []UserCommand
	commands = []UserCommand{
		/*
			{
				label: []string{"hud"},
				desc:  "Displays the HUD for the current game. Prints the game board and other game state information.",
				f: func() {

				}},
		*/

		{
			label:          []string{"settings"},
			desc:           "Displays Tourney's program settings.",
			classification: UI,
			f: func() {
				fmt.Print(Settings)
			}},
		{
			label: []string{"quit", "q", "exit"},
			desc:  "Stops any running tournament after the next game and quits the program.",
			f: func() {

				if T != nil && blocks(Tourneys.Selected().Done) {
					close(Tourneys.Selected().Done)
				}
				queQuit = true
			}},

		/*******************************************************************************

			.tourney Control Commands

		*******************************************************************************/

		{
			label:          []string{"info"},
			desc:           "Prints the configuration information for the selected tourney.",
			classification: TOURNEY_CONTROL,
			f: func() {
				Tourneys.Selected().Print()
			}},
		{
			label: []string{"play", "start"},
			desc:  "Starts playing the currently selected tourney on the local machine. Use the 'stop' command to stop playing.",
			f: func() {
				if len(words) > 1 {
					filename := words[1]
					Eval("load "+filename, Tourneys, wg)
				} else if T == nil {
					fmt.Println("No tournament is selected.")
					return
				}
				Tourneys.Selected().Done = make(chan struct{})
				wg.Add(1)
				go func() {
					if err := RunTourney(Tourneys.Selected()); err != nil {
						fmt.Println(err)
						fmt.Println("To add an engine to play in this tournament type: addengine")
					}
					wg.Done()
				}()
				return
			}},
		{
			label:          []string{"stop"},
			desc:           "Stops the tourney after the next game completes.",
			classification: TOURNEY_CONTROL,
			f: func() {
				if blocks(Tourneys.Selected().Done) {
					close(Tourneys.Selected().Done)
				}
			}},
		{
			label:          []string{"results", "standings"},
			desc:           "Displays the results of the currently selected tourney. Results may be displayed even if the tournament is incomplete.",
			classification: TOURNEY_CONTROL,
			f: func() {
				//standings := GenerateGameRecords(T, true)
				T.PlayerStandings.PrintStandings()
				//fmt.Print(standings.RenderTemplate())
				//fmt.Print(SummarizeResults(Tourneys.Selected()))
				fmt.Println("To see more details type, 'games' or 'g'")
			}},
		{
			label:          []string{"games"},
			desc:           "Displays the results of each game in the selected tourney.",
			classification: TOURNEY_CONTROL,
			f: func() {
				T.PrintGameList()
				//fmt.Print(SummarizeGames(Tourneys.Selected()))
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
			label:          []string{"loaddefault"},
			desc:           "Loads the default .tourney file.",
			f: func() {
				def, _ := LoadDefault()
				Tourneys.Add(def)
				ListActiveTourneys(Tourneys)
			}},
		{
			label:          []string{"save"},
			usage:          "save <filename>",
			example:        "save example.tourney",
			desc:           "Saves the currently selected tournaments settings as a .tourney file.",
			classification: TOURNEY_CONTROL,
			f: func() {
				fmt.Println("This command is not yet supported.")
				return

				var filename string
				if len(words) > 0 {
					filename = words[1]
				} else if Tourneys.Selected().filename != "" {
					filename = Tourneys.Selected().filename
				} else {
					fmt.Println("Please specify a filename. Type 'help save' for command usage info.")
				}
				Tourneys.Selected().filename = filename
				if err := SaveSettings(Tourneys.Selected()); err != nil {
					fmt.Println(err)
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
			label:          []string{"close"},
			desc:           "Closes the currently selected tournament.",
			classification: TOURNEY_CONTROL,
			f: func() {
				Eval("stop", Tourneys, wg)
				Tourneys.Remove()
				//fmt.Println("This command is not yet supported.")
			}},
		{
			label: []string{"ls"},
			desc:  "Displays a list of currently loaded tourneys.",
			f: func() {
				ListActiveTourneys(Tourneys)
			}},
		{
			label:          []string{"delete", "rm"},
			desc:           "Deletes all previous tourney data.",
			classification: TOURNEY_CONTROL,
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
			label:          []string{"engine", "add", "addengine"},
			desc:           "Adds an engine to the tournament.",
			classification: TOURNEY_CONTROL,
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
			label:          []string{"timecontrol", "time", "settime"},
			usage:          "timecontrol <moves>/<milliseconds>:<millisec added after each move>",
			example:        "time 40/2000:10",
			desc:           "Modifies the time control of the tournament. Time is measured in milliseconds.",
			classification: TOURNEY_CONTROL,
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
			label:          []string{"rounds"},
			usage:          "rounds <#>",
			example:        "rounds 10",
			desc:           "Modifies or Displays the number of rounds in the tournament.",
			classification: TOURNEY_CONTROL,
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

		/*******************************************************************************

			Help Commands:

		*******************************************************************************/

		{
			label: []string{"help"},
			desc:  "Displays help for a specific command.",
			f: func() {
				var str string
				if len(words) > 1 {
					var cmd UserCommand
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
					str = cmd.String()
				} else {
					str = "For a list of commands type 'commands'\nTo get help with a specific command type: help <command>\nFor even more help visit http://www.dirty-bit.com/tourney/support.html"
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

		/*******************************************************************************

			Network Commands:

		*******************************************************************************/

		{
			label: []string{"host"},
			desc:  "Hosts the currently loaded tournament on the local machine. However, does not start playing in the tournament.",
			f: func() {
				if len(words) > 1 {
					filename := words[1]
					Eval("load "+filename, Tourneys, wg)
				} else if Tourneys.Selected() == nil {
					fmt.Println("No tournament is selected.")
					return
				}
				wg.Add(1)
				go func() {
					Tourneys.Selected().Done = make(chan struct{})
					fmt.Print("\n\nTo stop hosting this tournament use the 'stop' command.\nThis machine can participate in the tournament by using the 'connect' command.")
					if err := HostTourney(Tourneys.Selected()); err != nil {
						fmt.Println(err)
					}
					wg.Done()
				}()
				return
			}},
		{
			label: []string{"connect", "c"},
			desc:  "Connects to a host running a tourney. Use the 'disconnect' command to disconnect.",
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
					Tourneys.ForceQuit = make(chan struct{})
					ConnectAndWait(address, Tourneys.ForceQuit)
					wg.Done()
				}()
				return
			}},
		{
			label: []string{"disconnect"},
			desc:  "Disconnect from a host.",
			f: func() {
				if blocks(Tourneys.ForceQuit) {
					close(Tourneys.ForceQuit)
				}
				return
			}},
		{
			label: []string{"dowork"},
			desc:  "Becomes a worker for the Dirt-Bit chess engine. Helps tune the engine.",
			f: func() {
				wg.Add(1)
				go func() {
					WorkForDirtyBit(Tourneys.ForceQuit)
					wg.Done()
				}()
				return
			}},
		{
			label: []string{"workers"},
			desc:  "Show the stats for connected workers.",
			f: func() {
				s := T.WorkerStats()
				for name, stats := range s {
					fmt.Print(name)
					if stats.Connected {
						fmt.Print("\t[Connected]\tAssigned Round ", stats.InProgress, "\t@ ", stats.Timestamp, "\n")
					} else {
						fmt.Print("\t[Disconnected]\n")
					}
				}
				if len(s) == 0 {
					fmt.Println("There are no workers for this tourney.")
				}
				return
			}},
			/*
		{
			label: []string{"broadcast"},
			desc:  "Broadcasts the currently selected tourney over http port " + strconv.Itoa(Settings.WebPort) + ". Broadcasting is enabled by default. The default port is specified in the 'tourney.settings' file.",
			f: func() {
				if !Tourneys.broadcasting {
					fmt.Println("Broadcasting http on port " + strconv.Itoa(Settings.WebPort))
					fmt.Println("Navigate your web browser to http://localhost:" + strconv.Itoa(Settings.WebPort))
					go func() {
						//if err := Broadcast(&T, selected); err != nil {
						if err := Broadcast(Tourneys); err != nil {
							fmt.Println(err)
						} else {

						}
						Tourneys.broadcasting = true
					}()

				} else {
					fmt.Println("Already broadcasting on port " + strconv.Itoa(Settings.WebPort))
				}
				return
			}},
			*/

		/*******************************************************************************

			Book Commands:

		*******************************************************************************/

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
			label: []string{"bookinfo"},
			desc:  "Shows information about the opening book for the current tourney.",
			f: func() {
				fmt.Println(T.OpeningBook.String())
				return
			}},
		/*******************************************************************************

			Testing only:

		*******************************************************************************/

		{
			label: []string{"test"},
			desc:  "Testbed for experimental functions. Used for developement only.",
			f: func() {
				T.PlayerStandings.PrintStandings()
				//R := GenerateGameRecords(T, true)
				//R.SortAllKeys()
				/*
					for k, v := range R.records {
						fmt.Println(k)
						fmt.Println("\t", v)
						fmt.Println()
					}

					fmt.Println("----------")
					standings := R.OverallStandings()
					fmt.Println(standings)
					fmt.Println()

					for _, engine := range T.Engines {
						player := engine.Name
						matchups := R.MatchupStandings(player)
						fmt.Println(player)
						fmt.Println(matchups)
						fmt.Println()
					}
					fmt.Println("----------")
				*/
				//fmt.Println(R.RenderTemplate())
				/*
					for player, _ := range R.records {
						for i, opponent := range R.orderedKeys[player] {
							rec := R.records[player][opponent]
							fmt.Println(i+1, player, "vs", opponent, "--> ", rec.Wins, "-", rec.Losses, "-", rec.Draws, "--> ", rec.Score(), "--> ", rec.Rate(), "%")
							fmt.Print(string(rec.Graph), "\n")
						}
						fmt.Println()
					}
				*/
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
			if command.classification == TOURNEY_CONTROL && T == nil {
				fmt.Println("No Tournament is selected.")
				fmt.Println("\tFor a list of commands type 'commands'.")
				fmt.Println("\tFor help with a command type 'help <command>'.")
				return queQuit
			}
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
	if len(Tourneys.List) > 0 {
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

	} else {
		fmt.Println("\nNo Tournaments are currently loaded.")
		fmt.Println("To load a .tourney file type 'load [filename]'.")
	}
	fmt.Println("For a list of additional commands, type 'commands'\n")

}

func PrintSysStats() {
	fmt.Print("System: ", runtime.GOMAXPROCS(0), "/", runtime.NumCPU(), " CPUs.")
	fmt.Println()
}
