/*

 Project: Tourney

 Module: moves
 Description: holds the move object and methods for interacting with it.
 	Eventually, engine data/logs will be tied into this?

 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*/

package main

type Move struct {
	algebraic string
	log       []string
}
