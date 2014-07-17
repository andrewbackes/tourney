/*

 Project: Tourney

 Module: engines
 Description: holds the engine object and methods for interacting with it
 
 Author(s): Andrew Backes, Daniel Sparks
 Created: 7/16/2014

*/

package main

type Engine struct {
	Name string
	Path string // file location
	Protocol string // uci or winboard
}

//Methods:
//	init() - go through all of the commands to get the engine ready to play
//	send() - send a command to the engine
//	wait() - wait for a certain response
//	idle() - idle loop. waits to be told to wait() or send().