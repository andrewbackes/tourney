/*******************************************************************************

 Project: Tourney

 Module: engines
 Description: Engine struct, Protocoler interface, and UCI/WinBoard structs.

 The Engine object has a Protocoler member. Then structs corresponding to UCI
 and Winboard impliment the Protocoler interface. The engine executable itself
 is ran in a goroutine, reader and writer Engine data members read/write the
 executables stdio so other Engine methods can interact with the executable.

TODO:
	-Error checking for it engine path exists and if it opens okay
	-WinBoard protocoler
	-Engines need to take options for hashtable size, multithreading, pondering,
	opening book, and a few other bare minimums.
	-Fix the bug where >> >> >> >> ... keeps looping sometimes.

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*******************************************************************************/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

/*******************************************************************************

	General Engine Functionality:

*******************************************************************************/

type Engine struct {
	//Public:
	Name     string
	Path     string // file location
	Protocol string // = "UCI" or "WINBOARD"

	//Private:
	reader *bufio.Reader
	writer *bufio.Writer
	logbuf *string

	protocol Protocoler // should = UCI{} or WINBOARD{}
	option   map[string]setting
}

type setting struct {
	optType    string
	optDefault string
	optMin     string
	optMax     string
}

func (E *Engine) Log(label string, record string) {
	*E.logbuf += fmt.Sprintln("[" + time.Now().Format("01/02/2006 15:04:05") + "][" + E.Name + "][" + label + "]" + record)
}

func (E *Engine) Evaluate(cmd string) error {
	E.Log("->", cmd)
	if cmd == "" {
		return nil
	}
	cmd = strings.ToLower(cmd)
	words := strings.Split(cmd, " ")
	switch words[0] {
	case "id":

	case "option":
		setupOption(E, words)

	case "info":

	}

	return nil
}

// when the engine says that it supports an option, add it to the struct:
func setupOption(E *Engine, words []string) {
	if n := getPair(words, "name"); n != "" {
		E.option[n] = setting{
			optDefault: getPair(words, "default"),
			optType:    getPair(words, "type"),
			optMin:     getPair(words, "min"),
			optMax:     getPair(words, "max")}
	}
}

func SetOption(E *Engine, name string, value string) error {
	if _, ok := E.option[strings.ToLower(name)]; ok {
		// TODO
	} else {
		// Does not support that option
		return errors.New("Not supported.")
	}
	return nil
}

// Helper :
func getPair(words []string, key string) string {
	for i, _ := range words {
		if words[i] == key && i+1 < len(words) {
			return words[i+1]
		}
	}
	return ""
}

// Send a command to the engine:
func (E *Engine) Send(s string) error {
	E.Log("<-", s)
	E.writer.WriteString(fmt.Sprintln(s)) // hopefully the line return is OS specific here.
	E.writer.Flush()
	//fmt.Print("->", fmt.Sprintln(s))
	return nil
}

// Recieve and process commands until a certain command is recieved
// or after the timeout (milliseconds) is achieved.
// Returns: last line read, lapsed time, error
func (E *Engine) Recieve(untilCmd string, timeout int64) (string, time.Duration, error) {

	var err error

	// Set up the timer:
	startTime := time.Now()
	// helper:
	type rec struct {
		data      string
		timestamp time.Time
	}

	// Start recieving from the engine:
	for {
		recieved := make(chan rec, 1)
		errChan := make(chan error, 1)
		/*
			//idea 1: (i think this has too much overhead)
				go func() {
					// TODO: determine how demanding this loop is on the system, try to minimize overhead.
					for {

						// check for something to read, so this goroutine doesnt just hault:
						// TODO: can still stall if there is something to read but it doesnt end with a '\n'
						_, e := E.reader.Peek(1) // < 1 byte to read => e != nil
						if e == nil {
							// Somthing to read:
							if nextline, err := E.reader.ReadString('\n'); err == nil {
								recieved <- nextline
								break
							} else {
								errChan <- err
								break
							}
						} else {
							// Nothing to read:
							select {
							case <-timedout:
								// stop reading
								break
							default:
								// keep reading
							}
						}
					}
					return
				}()
		*/

		//idea 2:
		go func() {
			// TODO: need a better idea here. ReadString() could hault this goroutine.
			if nextline, err := E.reader.ReadString('\n'); err == nil {
				recieved <- rec{nextline, time.Now()}
			} else {
				errChan <- err
			}
		}()

		/*
			//idea 3:
			go func() {
				buf := make([]byte, 1024)
				for {
					n, err := E.reader.Read(buf)
					if n != 0 {
						recieved <- string(buf[:n])
					}
					if err != nil {
						break
					}
				}
			}()
		*/
		//runtime.Gosched() // try to make sure the above goroutine gets priority

		// Since the timer and the reader are in goroutines, wait for:
		// (1) Something from the engine, (2) Too much time to pass. (3) An error
		select {
		case line := <-recieved:
			//l := lapsed()
			// Take off line return bytes:
			line.data = strings.Trim(line.data, "\r\n") // for windows
			line.data = strings.Trim(line.data, "\n")   // for *nix/bsd
			// Process the command recieved from the engine:
			if err = E.Evaluate(line.data); err != nil {
				return "", line.timestamp.Sub(startTime), errors.New("Error recieving from engine: " + err.Error())
			}
			// Check if the recieved command is the one we are waiting for:
			if strings.Contains(line.data, untilCmd) {
				return line.data, line.timestamp.Sub(startTime), nil
			}
		case <-time.After(time.Duration(timeout) * time.Millisecond):
			return "", time.Now().Sub(startTime), errors.New("Timed out waiting for engine to respond.")
		case e := <-errChan:
			return "", time.Now().Sub(startTime), errors.New("Error recieving from engine: " + e.Error())
		}
	}
	return "", time.Now().Sub(startTime), nil
}

// Set the engine up to be ready to think on its first move:
func (E *Engine) Start(logbuffer *string) error {
	E.logbuf = logbuffer

	E.option = make(map[string]setting)

	// Decide which protocol to use:
	// TODO: add some autodetect code here
	if strings.ToUpper(E.Protocol) == "UCI" {
		E.protocol = UCI{}
	} else if strings.ToUpper(E.Protocol) == "WINBOARD" {
		//E.protocol = WINBOARD{}
	}

	cmd := exec.Command(E.Path)

	// Setup the pipes to communicate with the engine:
	StdinPipe, errIn := cmd.StdinPipe()
	if errIn != nil {
		return errors.New("Error Initializing Engine: can not establish in pipe.")
	}
	StdoutPipe, errOut := cmd.StdoutPipe()
	if errOut != nil {
		return errors.New("Error Initializing Engine: can not establish out pipe.")
	}
	E.writer, E.reader = bufio.NewWriter(StdinPipe), bufio.NewReader(StdoutPipe)

	// Start the engine:
	started := make(chan struct{})
	errChan := make(chan error)
	go func() {
		// Question: Does this force the engine to run in its own thread?
		if err := cmd.Start(); err != nil {
			errChan <- err
			return
			//return errors.New("Error executing " + E.Path + " - " + err.Error())
		}
		close(started)
	}()
	select {
	case <-started:
	case e := <-errChan:
		return errors.New("Error executing " + E.Path + " - " + e.Error())
	}

	// Get the engine ready:
	s, r := E.protocol.Initialize()
	E.Send(s)
	E.Recieve(r, 2500)
	E.NewGame()

	// Setup up for when the engine exits:
	go func() {
		cmd.Wait()
		//TODO: add some confirmation that the engine has terminated correctly.
	}()

	return nil
}

func (E *Engine) NewGame() error {
	//E.protocol.New(E.reader, E.writer)
	E.Send(E.protocol.NewGame())
	return nil
}

// The engine should close itself:
func (E *Engine) Shutdown() error {
	// TODO: add confirmation that the engine has shut down correctly
	E.Send(E.protocol.Quit())
	return nil
}

// The engine should decide what move it wants to make:
func (E *Engine) Move(timers [2]int64, movesToGo int64) (Move, time.Duration, error) {
	s, r := E.protocol.Move(timers, movesToGo)
	E.Send(s)
	max := timers[WHITE]
	if timers[BLACK] > max {
		max = timers[BLACK]
	}
	response, t, err := E.Recieve(r, max+1000)
	if err != nil {
		return Move{}, t, err
	}

	// figure out what move was picked:
	words := strings.Split(response, " ") // bestmove e2e4 ponder e7e5
	chosenMove := Move{Algebraic: words[1]}
	return chosenMove, t, nil
}

// The engine should set its internal board to adjust for the moves far in the game
func (E *Engine) Set(movesSoFar []Move) error {
	s := E.protocol.SetBoard(movesSoFar)
	err := E.Send(s)
	return err
}

func (E *Engine) Ping() error {
	s, r := E.protocol.Ping()
	E.Send(s)
	_, _, err := E.Recieve(r, -1)
	return err
}

/*******************************************************************************

	Protocol Specific:

*******************************************************************************/

type Protocoler interface {
	Initialize() (string, string)
	Quit() string
	Move(timers [2]int64, movesToGo int64) (string, string)
	SetBoard(moveSoFar []Move) string
	NewGame() string
	Ping() (string, string)
}

type UCI struct{}
type WINBOARD struct{}

/*******************************************************************************

	UCI:

*******************************************************************************/

func (U UCI) Ping() (string, string) {
	return "isready", "readyok"
}

func (U UCI) Initialize() (string, string) {
	// (command to send),(command to recieve)
	return "uci", "uciok"
}

func (U UCI) NewGame() string {
	return "ucinewgame"
}

func (U UCI) Quit() string {
	return "quit"
}

func (U UCI) Move(timer [2]int64, movesToGo int64) (string, string) {
	goString := "go"

	if timer[WHITE] > 0 {
		goString += " wtime " + strconv.FormatInt(timer[WHITE], 10)
	}
	if timer[BLACK] > 0 {
		goString += " btime " + strconv.FormatInt(timer[BLACK], 10)
	}
	if movesToGo > 0 {
		goString += " movestogo " + strconv.FormatInt(movesToGo, 10)
	}

	return goString, "bestmove"
}

func (U UCI) SetBoard(movesSoFar []Move) string {
	var ml []string

	for _, m := range movesSoFar {
		ml = append(ml, m.Algebraic)
	}

	var pos string
	if len(movesSoFar) > 0 {
		pos = "position startpos moves " + strings.Join(ml, " ")
	} else {
		pos = "position startpos"
	}

	return pos

}

/*******************************************************************************

	WINBOARD:

*******************************************************************************/

/* TODO: EVERYTHING */
