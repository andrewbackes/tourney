/*******************************************************************************

 Project: Tourney

 Module: winboard
 Description: implements Protocoler from engines.go

TODO:


 Author(s): Andrew Backes
 Created: 12/15/2014

*******************************************************************************/

package main

import (
	//"bufio"
	//"errors"
	//"fmt"
	//"os"
	//"os/exec"
	"strconv"
	"strings"
	//"time"
)

/*******************************************************************************

	WINBOARD:

*******************************************************************************/

type WINBOARD struct {
	features map[string]string
}

func (W WINBOARD) Initialize() (string, func(parse string) bool) {
	recieve := "done=1"
	return "xboard\nprotover 2", func(parse string) bool {
		return len(parse) >= len(recieve) && parse[len(parse)-len(recieve):] == recieve
	}
}

func (W *WINBOARD) Move(Timer [2]int64, MovesToGo int64, EngineColor Color) (string, func(parse string) bool) {
	var goString string

	// Note: remember that winboard uses centiseconds
	goString += "time " + strconv.FormatInt(Timer[EngineColor]/10, 10) + "\n"
	goString += "otim " + strconv.FormatInt(Timer[[]int{1, 0}[EngineColor]]/10, 10) + "\n"
	if W.features["colors"] == "1" {
		goString += []string{"white\n", "black\n"}[EngineColor]
	}
	goString += "go"

	recieve := "move "
	return goString, func(parse string) bool {
		moved := len(parse) >= len(recieve) && parse[:len(recieve)] == recieve
		resigned := strings.Contains(parse, "resign")
		return moved || resigned
	}
}

func (W WINBOARD) Ping(N int) (string, func(parse string) bool) {
	if W.features["ping"] == "1" {
		recieve := "pong " + strconv.Itoa(N)
		return "ping " + strconv.Itoa(N), func(parse string) bool {
			return len(parse) >= len(recieve) && parse[:len(recieve)] == recieve
		}
	}
	return "", func(parse string) bool { return false }
}

func (W WINBOARD) NewGame(Time, Moves int64) string {
	level := "level " + strconv.FormatInt(Moves, 10) + " " + strconv.FormatInt((Time/1000)/60, 10) + ":" + strconv.FormatInt((Time/1000)%60, 10) + " 0\n"
	return "new\nrandom\n" + level + "post\nhard\neasy\ncomputer"
}

func (W *WINBOARD) SetBoard(movesSoFar []Move, analysisSoFar []MoveAnalysis) string {
	// TODO: instead of passing []MoveAnalysis say how many moves in the book there are.

	var pos string

	// Determine if this is the first move this engine will be thinking on:
	movesOutOfBook := 0
	pos = "force\n"
	for i, _ := range movesSoFar {
		if v := W.features["usermove"]; v == "1" {
			pos += "usermove "
		}
		pos += string(movesSoFar[i]) + "\n"
		if analysisSoFar[i].Comment != BOOKMOVE {
			movesOutOfBook++
		}
	}
	pos = strings.TrimSuffix(pos, "\n")

	// when there is more than one move out of the book, dont play the opening:
	if movesOutOfBook > 1 {
		pos = "force\n"
		if v := W.features["usermove"]; v == "1" {
			pos += "usermove "
		}
		pos += string(movesSoFar[len(movesSoFar)-1]) //only the last move is needed
	}

	return pos
}

func (W WINBOARD) Quit() string {
	return "quit"
}

func parsePost(line string) EvaluationData {
	// ply score time nodes pv
	// ex: 6    11       0      5118 Qd5 9.Bf4 Nc6 10.e3 Bg4 11.a3 [TT]
	// ex: 8&     66    1    20536   d1e2  e8e7  e2e3  e7e6  e3d4  g7g5  a2a4  f7f5

	fields := strings.Fields(line)

	var d, s, t, n int
	var l, u bool
	var pv string
	var err error

	if len(fields) >= 4 {
		if d, err = strconv.Atoi(fields[0]); err != nil {
			d, _ = strconv.Atoi(fields[0][:len(fields[0])-1])
		}
		s, _ = strconv.Atoi(fields[1])
		t, _ = strconv.Atoi(fields[2])
		n, _ = strconv.Atoi(fields[3])
	}
	for i := 4; i < len(fields); i++ {
		pv += fields[i] + " "
	}
	pv = strings.Trim(pv, " ")
	if len(pv) > 1 {
		l = (pv[len(pv)-1] == '!')
		u = (pv[len(pv)-1] == '?')
	}

	return EvaluationData{D: d, V: s, U: u, L: l, T: t, N: n, P: pv}
}

func (W WINBOARD) ExtractMove(output string) (Move, MoveAnalysis) {

	output = strings.Replace(output, "\r", " ", -1)
	lines := strings.Split(output, "\n")
	var mv Move
	var ma MoveAnalysis
	for _, line := range lines {
		if strings.HasPrefix(line, "move ") {
			words := strings.Fields(line)
			//words := strings.Split(line, " ")
			if len(words) >= 2 {
				mv = Move(words[1])
			}
			break
		} else if strings.Contains(line, "resign") {
			mv = "0000"
		} else {
			eval := parsePost(line)
			if eval.Depth() != 0 {
				ma.Evaluation = append(ma.Evaluation, eval)
			}
		}
	}
	return mv, ma
}

func (W *WINBOARD) RegisterEngineOptions(output string, options map[string]Setting) {

	// helper. Splits based on spaces not inside quotes:
	nonQuotedWordSplit := func(ln string) []string {
		r := []string{}
		quoted := false
		var b int
		for i, v := range ln {
			if string(v) == "\"" {
				quoted = !quoted
			}
			if string(v) == " " && !quoted || i == len(ln)-1 {
				r = append(r, strings.Trim(ln[b:i+1], " "))
				b = i + 1
			}
		}
		return r
	}
	// ***

	W.features = make(map[string]string) // init for local struct use
	W.setFeaturesToDefault()

	output = strings.Replace(output, "\r", "", -1)
	lines := strings.Split(output, "\n")

	for _, v := range lines {
		if strings.HasPrefix(v, "feature") {
			v = v[len("feature "):]
		} else {
			continue
		}

		pairs := nonQuotedWordSplit(v)
		for j, _ := range pairs {
			p := strings.Split(pairs[j], "=")
			if p[0] != "option" { // TEMPORARY
				if len(p) > 1 {
					W.features[p[0]] = p[1]
					//fmt.Println("accepted", p[0], p[1])
				}
			}
		}
	}
}

// Sets the feature list to the Winboard defaults:
func (W *WINBOARD) setFeaturesToDefault() {

	// Winboard/xboard default values:
	W.features["ping"] = "0"      //ping (boolean, default 0, recommended 1)
	W.features["setboard"] = "0"  //setboard (boolean, default 0, recommended 1)
	W.features["playother"] = "0" //playother (boolean, default 0, recommended 1)
	W.features["san"] = "0"       //san (boolean, default 0)
	W.features["usermove"] = "0"  //usermove (boolean, default 0)
	W.features["time"] = "1"      //time (boolean, default 1, recommended 1)
	W.features["draw"] = "1"      //draw (boolean, default 1, recommended 1)
	W.features["sigint"] = "1"    //sigint (boolean, default 1)
	W.features["sigterm"] = "1"   //sigterm (boolean, default 1)
	W.features["reuse"] = "1"     //reuse (boolean, default 1, recommended 1)
	W.features["analyze"] = "1"   //analyze (boolean, default 1, recommended 1)
	W.features["colors"] = "1"    //colors (boolean, default 1, recommended 0)
	W.features["ics"] = "0"       //ics (boolean, default 0)
	W.features["pause"] = "0"     // pause (boolean, default 0)
	W.features["debug"] = "0"     //debug (boolean, default 0)
	W.features["memory"] = "0"    //memory (boolean, default 0)
	W.features["smp"] = "0"       //smp (boolean, default 0)
	W.features["exclude"] = "0"   //exclude (boolean, default 0)
	W.features["setscore"] = "0"  //setscore (boolean, default 0)
	W.features["highlight"] = "0" //highlight (boolean, default 0)

}
