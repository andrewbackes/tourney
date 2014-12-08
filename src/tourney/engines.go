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
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// helper:
type rec struct {
	timestamp time.Time
	data      string
}

/*******************************************************************************

	General Engine Functionality:

*******************************************************************************/

type Engine struct {
	//Public:
	Name     string
	Path     string // file location
	Protocol string // = "UCI" or "WINBOARD"
	MD5      string

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

func (E *Engine) ValidateEngineFile() error {
	// First decides if the file exists.
	// Compares the checksum to the md5 sum that is stored in memory.
	// If nothing has been stored, then it saves this checksum.
	// Returns true when they match or it was previously blank.

	// Existence:
	if _, err := os.Stat(E.Path); os.IsNotExist(err) {
		return err
	}

	// Check sum:
	if checksum, err := GetMD5(E.Path); err != nil {
		return err
	} else {
		if E.MD5 == "" {
			// md5 has not been previously checked
			E.MD5 = checksum
		} else if E.MD5 != checksum {
			return errors.New("MD5 mismatch")
		}
	}

	return nil
}

func (E *Engine) Log(label string, record rec) {
	*E.logbuf += fmt.Sprintln("[" + record.timestamp.Format("01/02/2006 15:04:05.000") + "][" + E.Name + "][" + label + "]" + record.data)
}

// INCOMPLETE:
func (E *Engine) Evaluate(cmd string) error {
	//E.Log("->", rec{time.Now(), cmd})
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
	E.Log("<-", rec{time.Now(), s})
	E.writer.WriteString(fmt.Sprintln(s)) // hopefully the line return is OS specific here.
	E.writer.Flush()
	//fmt.Print("->", fmt.Sprintln(s))
	return nil
}

// Recieve and process commands until a certain command is recieved
// or after the timeout (milliseconds) is achieved.
// Returns: engine output, lapsed Time, error
func (E *Engine) Recieve(untilCmd string, timeout int64) (string, time.Duration, error) {

	//var err error
	var output string //engine's output

	// Set up the Timer:
	startTime := time.Now()

	// Start recieving from the engine:
	for {
		recieved := make(chan rec, 1)
		errChan := make(chan error, 1)

		//TODO: Redesign neccessary
		go func() {
			// TODO: need a better idea here. ReadString() could hault this goroutine.
			if nextline, err := E.reader.ReadString('\n'); err == nil {
				recieved <- rec{time.Now(), nextline}
			} else {
				errChan <- err
			}
		}()

		// Since the Timer and the reader are in goroutines, wait for:
		// (1) Something from the engine, (2) Too much Time to pass. (3) An error
		select {
		case line := <-recieved:
			// keep track of the total output from the engine:
			output += line.data

			// Take off line return bytes:
			line.data = strings.Trim(line.data, "\r\n") // for windows
			line.data = strings.Trim(line.data, "\n")   // for *nix/bsd

			// Log this line of engine output:
			E.Log("->", line)

			// Check if the recieved command is the one we are waiting for:
			if strings.Contains(line.data, untilCmd) {
				return output, line.timestamp.Sub(startTime), nil
			}

		case <-time.After(time.Duration(timeout) * time.Millisecond):
			description := "Timed out waiting for engine to respond."
			E.Log("ERROR", rec{time.Now(), description})
			return output, time.Now().Sub(startTime), errors.New(description)

		case e := <-errChan:
			description := "Error recieving from engine: " + e.Error()
			E.Log("ERROR", rec{time.Now(), description})
			return output, time.Now().Sub(startTime), errors.New(description)

		}
	}
	return output, time.Now().Sub(startTime), nil //this should never occur
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
	// TODO: evaluate output. This should set up the options for the engine.
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
func (E *Engine) Move(timers [2]int64, MovesToGo int64) (Move, time.Duration, error) {
	s, r := E.protocol.Move(timers, MovesToGo)
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
	return E.protocol.ExtractMove(response), t, nil
}

// The engine should set its internal Board to adjust for the Moves far in the game
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
	Move(timers [2]int64, MovesToGo int64) (string, string)
	SetBoard(moveSoFar []Move) string
	NewGame() string
	Ping() (string, string)

	ExtractMove(string) Move
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

func (U UCI) Move(Timer [2]int64, MovesToGo int64) (string, string) {
	goString := "go"

	if Timer[WHITE] > 0 {
		goString += " wtime " + strconv.FormatInt(Timer[WHITE], 10)
	}
	if Timer[BLACK] > 0 {
		goString += " btime " + strconv.FormatInt(Timer[BLACK], 10)
	}
	if MovesToGo > 0 {
		goString += " movestogo " + strconv.FormatInt(MovesToGo, 10)
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

func (U UCI) ExtractMove(output string) Move {
	// figure out what move was picked:
	//lines := strings.Split(response, "\n")
	//words := strings.Split(lines[len(lines)-2], " ") // bestmove e2e4 ponder e7e5
	//chosenMove := Move{Algebraic: words[1]}

	// TODO: REFACTOR: this replace also happens in Engine.Recieve()
	output = strings.Replace(output, "\n\r", " ", -1)
	output = strings.Replace(output, "\n", " ", -1)
	words := strings.Split(output, " ")

	// Helper function:
	LastValueOf := func(key string) string {
		//returns the word after the word given as an arg
		for i := len(words) - 1; i >= 0; i-- {
			if words[i] == key {
				if i+1 <= len(words)-1 {
					return words[i+1]
				}
			}
		}
		return ""
	}

	d, _ := strconv.Atoi(LastValueOf("depth"))
	t, _ := strconv.Atoi(LastValueOf("time")) // TODO: doing this way may only give the time for this depth
	skey := LastValueOf("score")              // ex: score cp 112   but it could be:   score mate 7

	var sval int
	if skey == "cp" {
		sval, _ = strconv.Atoi(LastValueOf(skey))
	} else if skey == "mate" {
		sval, _ = strconv.Atoi(LastValueOf(skey))
		sval = MateIn(sval)
	}

	return (Move{
		Algebraic: LastValueOf("bestmove"),
		Depth:     d,
		Time:      t,
		Score:     sval,
	})

}

/*******************************************************************************

	WINBOARD:

*******************************************************************************/

/* TODO: EVERYTHING */
