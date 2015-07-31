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
	"net/url"
	"path/filepath"
	"strconv"
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
	tmpl, err := template.ParseFiles(page, filepath.Join(Settings.TemplateDirectory, "_header.html"))
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

func requestHandler(w http.ResponseWriter, req *http.Request, t *Tourney) {
	q, err := url.ParseQuery(req.URL.RawQuery)
	if err != nil {
		//log.Println("Request Error: ", req.RemoteAddr, err, req.URL.RawQuery)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	var query = make(map[string]string)
	for k, v := range q {
		query[k] = v[0]
	}
	round, _ := strconv.Atoi(query["round"])
	ply, _ := strconv.Atoi(query["ply"])
	if round == 0 {
		round = 1
	}
	if ply == 0 {
		ply = 1
	}
	if query["display"] == "standings" || query["display"] == "" {
		renderTourneyPage(w, t)
	} else if query["display"] == "game" {
		renderGameViewer(w, t, round)
	} else if query["display"] == "round" {
		renderRoundPage(w, t, round)
	} else if query["display"] == "ply" {
		renderPlyPage(w, t, round, ply)
	} else if query["display"] == "gamelist" {
		renderGameListPage(w, t)
	}

}

// Broadcast turns starts serving http for the tourney data.
// examples:
// 		http://localhost/view?display=standings
// 		http://localhost/view?display=round&round=1
// 		http://localhost/view?display=ply&ply=1
// 		http://localhost/view?display=game&round=1
// 		http://localhost/view?display=log&round=1
func Broadcast(Tourneys *TourneyList) error {
	//TODO: check that the tourney is valid

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/view", http.StatusFound)
	})

	http.HandleFunc("/view", func(w http.ResponseWriter, req *http.Request) {
		requestHandler(w, req, Tourneys.Selected())
	})

	// Set up a file server for resources such as scripts, images, etc.
	http.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir(filepath.Join(Settings.TemplateDirectory, "resources")))))
	// Log Requests:
	http.Handle("/logs/", http.StripPrefix("/logs/", http.FileServer(http.Dir(Settings.LogDirectory))))

	// Start the server:
	// TODO: allow the server to be shut down.
	err := http.ListenAndServe(":"+strconv.Itoa(Settings.WebPort), nil)
	if err != nil {
		return err
	}

	return nil
}
