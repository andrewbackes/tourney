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
 	-a search log file for move function.

 BUGS:
 	-executing Broadcast() more than once crashes.

*******************************************************************************/

package main

import (
	"fmt"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func renderNothingLoaded(w http.ResponseWriter) {
	io.WriteString(w, "No Tournament is loaded.\n")
	io.WriteString(w, "From within Tourney, you can type: 'load [filename]' to load a .tourney file.\n")
	io.WriteString(w, "To create a new .tourney file type 'new'.\n")
	io.WriteString(w, "\n")
	io.WriteString(w, "For a complete list of commands type 'commands'.\n")
	io.WriteString(w, "Get help with a command by typing 'help [command]'.\n")
}

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
	if T == nil {
		renderNothingLoaded(w)
		return
	}
	renderTemplate(w, filepath.Join(Settings.TemplateDirectory, "tourney.html"), T)

}

func renderRoundPage(w http.ResponseWriter, T *Tourney, round int) {
	if T == nil {
		renderNothingLoaded(w)
		return
	}
	if round < len(T.GameList) && round >= 0 {
		renderTemplate(w, filepath.Join(Settings.TemplateDirectory, "game.html"), T.GameList[round-1])
	} else {
		io.WriteString(w, "That is not a valid round in this Tourney.")
	}
}

func renderGameListPage(w http.ResponseWriter, T *Tourney) {
	if T == nil {
		renderNothingLoaded(w)
		return
	}
	renderTemplate(w, filepath.Join(Settings.TemplateDirectory, "gamelist.html"), T)
}

func renderPlyPage(w http.ResponseWriter, T *Tourney, round, ply int) {
	if T == nil {
		renderNothingLoaded(w)
		return
	}
	if round < len(T.GameList) && round >= 0 {
		renderTemplate(w, filepath.Join(Settings.TemplateDirectory, "ply.html"), T.GameList[round-1].AnalysisList[ply])
	} else {
		io.WriteString(w, "That is not a valid ply of a round in this Tourney.")
	}
}

func renderGameViewer(w http.ResponseWriter, T *Tourney, round int) {
	if T == nil {
		renderNothingLoaded(w)
		return
	}
	if round < len(T.GameList) && round >= 0 {
		renderTemplate(w, filepath.Join(Settings.TemplateDirectory, "viewer.html"), T.GameList[round-1])
	} else {
		io.WriteString(w, "That is not a valid round in this Tourney.")
	}
}

//func Broadcast(T *Tourney) error {
//func Broadcast(TList *[]*Tourney, Tindex *int) error {
func Broadcast(Tourneys *TourneyList) error {
	//TODO: check that the tourney is valid

	// Summary Requests:
	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		renderTourneyPage(w, Tourneys.Selected())
	})

	// Game History Requests:
	http.HandleFunc("/gamelist/", func(w http.ResponseWriter, req *http.Request) {
		renderGameListPage(w, Tourneys.Selected())
	})

	// Round Requests:
	http.HandleFunc("/round/", func(w http.ResponseWriter, req *http.Request) {
		request := strings.Trim(req.URL.Path[len("/round"):], "/")
		words := strings.Split(request, "/")
		if len(words) == 1 {
			// just the round is being requested:
			round, _ := strconv.Atoi(words[0])
			renderRoundPage(w, Tourneys.Selected(), round)
		} else if len(words) >= 3 && words[1] == "ply" {
			round, _ := strconv.Atoi(words[0])
			ply, _ := strconv.Atoi(words[2])
			renderPlyPage(w, Tourneys.Selected(), round, ply)
		}
	})

	// Game Viewer:
	http.HandleFunc("/viewer/", func(w http.ResponseWriter, req *http.Request) {
		request, _ := strconv.Atoi(strings.Trim(req.URL.Path[len("/viewer"):], "/"))
		renderGameViewer(w, Tourneys.Selected(), request)
	})
	// Image files for Game Viewer:
	http.Handle("/viewer/pieces/", http.StripPrefix("/viewer/pieces/", http.FileServer(http.Dir(filepath.Join(Settings.TemplateDirectory, "pieces")))))

	// Log Requests:
	//http.Handle("/logs/", http.FileServer(http.Dir("./")))
	http.Handle("/logs/", http.StripPrefix("/logs/", http.FileServer(http.Dir(Settings.LogDirectory))))

	// Start the server:
	// TODO: allow the server to be shut down.
	err := http.ListenAndServe(":"+strconv.Itoa(Settings.WebPort), nil)
	if err != nil {
		return err
	}

	return nil
}
