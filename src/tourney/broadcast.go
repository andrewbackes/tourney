/*******************************************************************************

 Project: Tourney
 Author(s): Andrew Backes
 Created: 11/22/2014

 Module: broadcast
 Description: Web broadcasting services.

 TODO:
	-need to be able to disable the web server from within the program.
 	-Push move by move in games (Server-Sent).
 	-Have a 'compare pv' view where the moves are lined up on top of eachother.

*******************************************************************************/

package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func renderTemplate(w http.ResponseWriter, page string, obj interface{}) {
	tmpl, err := template.ParseFiles(page)
	if err != nil {
		fmt.Println(err)
		io.WriteString(w, fmt.Sprint("Error opening '", page, "' - ", err))
		return
	}
	err = tmpl.Execute(w, obj)
	if err != nil {
		fmt.Println(err)
		io.WriteString(w, fmt.Sprint("Error executing parse on '", page, "' - ", err))
		return
	}
}

func renderTourneyPage(w http.ResponseWriter, T *Tourney) {
	//io.WriteString(w, SummarizeResults(T))
	//io.WriteString(w, SummarizeGames(T))
	Records := NewRecordRollup(T)
	renderTemplate(w, "templates/tourney.html", Records)
}

func renderRoundPage(w http.ResponseWriter, T *Tourney, round int) {
	if round < len(T.GameList) && round >= 0 {
		//io.WriteString(w, "Round: "+strconv.Itoa(round))
		//io.WriteString(w, fmt.Sprint(T.GameList[round].MoveList))
		renderTemplate(w, "templates/game.html", T.GameList[round-1])
		//renderTemplate(w, "templates/viewer.html", T.GameList[round])
	} else {
		io.WriteString(w, "That is not a valid round in this Tourney.")
	}
}

func renderGameViewer(w http.ResponseWriter, T *Tourney, round int) {
	if round < len(T.GameList) && round >= 0 {
		renderTemplate(w, "templates/viewer.html", T.GameList[round-1])
	} else {
		io.WriteString(w, "That is not a valid round in this Tourney.")
	}
}

//func Broadcast(T *Tourney) error {
func Broadcast(TList []*Tourney, Tindex *int) error {
	//TODO: check that the tourney is valid

	// Summary Requests:
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		renderTourneyPage(w, TList[*Tindex])
	})

	// Round Requests:
	http.HandleFunc("/round/", func(w http.ResponseWriter, req *http.Request) {
		request, _ := strconv.Atoi(strings.Trim(req.URL.Path[len("/round"):], "/"))
		renderRoundPage(w, TList[*Tindex], request)
	})

	// Game Viewer:
	http.HandleFunc("/viewer/", func(w http.ResponseWriter, req *http.Request) {
		request, _ := strconv.Atoi(strings.Trim(req.URL.Path[len("/viewer"):], "/"))
		renderGameViewer(w, TList[*Tindex], request)
	})
	// Image files for Game Viewer:
	http.Handle("/viewer/pieces/", http.StripPrefix("/viewer/pieces/", http.FileServer(http.Dir("templates/pieces"))))

	// Log Requests:
	http.Handle("/logs/", http.FileServer(http.Dir("./")))

	// Start the server:
	// TODO: allow the server to be shut down.
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		return err
	}

	return nil
}
