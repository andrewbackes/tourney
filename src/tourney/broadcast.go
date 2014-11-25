/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes
 Created: 11/22/2014

 Module: broadcast
 Description: Web broadcasting services.

 TODO:
 	-Push move by move in games (Server-Sent).

*******************************************************************************/

package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func renderTourneyPage(w http.ResponseWriter, T *Tourney) {
	io.WriteString(w, SummarizeResults(T))
	io.WriteString(w, SummarizeGames(T))
}

func renderRoundPage(w http.ResponseWriter, T *Tourney, round int) {

	if round < len(T.GameList) && round >= 0 {
		io.WriteString(w, "Round: "+strconv.Itoa(round))
		io.WriteString(w, fmt.Sprint(T.GameList[round].MoveList))
	} else {
		io.WriteString(w, "That is not a valid round in this Tourney.")
	}
}

func Broadcast(T *Tourney) error {
	//TODO: check that the tourney is valid

	// Summary Requests:
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		renderTourneyPage(w, T)
	})

	// Round Requests:
	http.HandleFunc("/round/", func(w http.ResponseWriter, req *http.Request) {
		request, _ := strconv.Atoi(strings.Trim(req.URL.Path[len("/round"):], "/"))
		renderRoundPage(w, T, request)
	})

	// Start the server:
	// TODO: allow the server to be shut down.
	err := http.ListenAndServe(":8000", nil)
	if err != nil {
		return err
	}

	return nil
}
