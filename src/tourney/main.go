// Project: Tourney
//
// Author(s): Andrew Backes, Daniel Sparks
// Created: 7/14/2014

package main

import (
	"fmt"
)

func main() {

	fmt.Println("Project: Tourney Started")

	// Until there is a need to have multiple Tourney objects to run at once,
	// this single object will just be passed around and manipulated:
	var tourney Tourney

	// Check for a lanuch arguement with for a .tourney file
	// .tourney files contain all of the settings needed
	// to start a tourney without any terminal input.

	// validate that the file exists and is valid:

	// when no .tourney file is provided or is invalid, should load default.tourney
	tourney.LoadDefault()

	// TODO: Other launch arguements

	// and either go to the menu or the command loop
	commandLoop(&tourney)

}

func printError(description string, err string) {
	fmt.Println("ERROR:", err, "-", description)
}
