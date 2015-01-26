/*******************************************************************************

 Project: Tourney

 Module: evaluation
 Description: storage of engine's evaluation data

TODO:

 Author(s): Andrew Backes
 Created: 01/25/2015

*******************************************************************************/

package main

/*
import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/exec"
	//"strconv"
	"strings"
	"time"
)
*/

//
// The data members are non descriptive to save space when encoding to json.
// This is important not only for disk space but also because these objects
// are sent back and forth in network tournaments.
//
type EvaluationData struct {
	// Depth:
	D int
	// SelDepth:
	S int `json:",omitempty"`
	// Score:
	V int
	// Lowerbound:
	L bool `json:",omitempty"`
	// Upperbound:
	U bool `json:",omitempty"`
	// Time:
	T int
	// Nodes:
	N int `json:",omitempty"`
	// PV:
	P string `json:",omitempty"`
}

//
// Since the data members are not descriptive, the following methods can be used
//
func (E *EvaluationData) Depth() int {
	return E.D
}
func (E *EvaluationData) SelDepth() int {
	return E.S
}
func (E *EvaluationData) Seldepth() int {
	return E.S
}
func (E *EvaluationData) Score() int {
	return E.V
}
func (E *EvaluationData) Lowerbound() bool {
	return E.L
}
func (E *EvaluationData) Upperbound() bool {
	return E.U
}
func (E *EvaluationData) Time() int {
	return E.T
}
func (E *EvaluationData) Nodes() int {
	return E.N
}
func (E *EvaluationData) Pv() string {
	return E.P
}
