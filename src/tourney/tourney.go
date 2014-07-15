// Project: Tourney
//
// Author(s): Andrew Backes, Daniel Sparks
// Created: 7/14/2014


package main

import "fmt"
import "bitboard"

func main() {
	var bb uint64 = 1337
	fmt.Println(bb)
	fmt.Println("Project Tourney Started")
	bitboard.Print(bb)
}
