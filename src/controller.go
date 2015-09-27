/*******************************************************************************

 Project: 		Tourney
 Module: 		Controller
 Created: 		9/27/2015
 Author(s): 	Andrew Backes
 Description: 	Interfaces between the web ui / console ui and the data model.

*******************************************************************************/

package main

import (
	"sync"
)

type Controller struct {
	commandQue chan string
	tourneys   TourneyList
	wg         sync.WaitGroup
	quit       chan struct{}
}

func NewController() Controller {
	return Controller{
		commandQue: make(chan string , 16),
		quit:       make(chan struct{}),
	}
}

func (c *Controller) Enque(command string) {
	c.commandQue <- command
}

func (c *Controller) Start() {
	quit := make(chan struct{})
	done := false
	for !done {
		select {
		case <-quit:
		case command := <-c.commandQue:
			done = Eval(command, &c.tourneys, &c.wg)
			if done {
				close(quit)
			}
		}
	}
	c.wg.Wait()
}

func (c *Controller) Stopped() bool {
	return !blocks(c.quit)
}

func (c *Controller) PromptString() string {
	if t := c.tourneys.Selected(); t != nil {
		return t.Event + "> "
	}
	return "> "
}
