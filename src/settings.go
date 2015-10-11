/*******************************************************************************

 Project: Tourney
 Module: games
 Author(s): Andrew Backes
 Created: 1/7/2015

 Description: Global settings control.

*******************************************************************************/

package main

import (
	"strconv"
	"strings"
	"encoding/json"
	"fmt"
	"os"
)

type GlobalSettings struct {
	WorkerDirectory   string
	LogDirectory      string
	TemplateDirectory string
	TourneyDirectory  string
	SaveDirectory     string
	BookDirectory     string
	BuildDirectory	  string
	
	ServerPort     int
	WebPort        int
	EngineFilePort int

	MaxConnectionAttempts int
}

func DefaultSettings() GlobalSettings {
	return GlobalSettings{
		WorkerDirectory:       "worker/",
		LogDirectory:          "logs/",
		TemplateDirectory:     "templates/",
		TourneyDirectory:      "tourneys/",
		SaveDirectory:         "data/",
		BookDirectory:         "book/",
		BuildDirectory:        "build/",
		ServerPort:            9000,
		WebPort:               8080,
		EngineFilePort:        9001,
		MaxConnectionAttempts: 3,
	}
}

//
// String() is used to print with the fmt package
//
func (G GlobalSettings) String() string {
	title := " Tourney's Program Settings: "
	return strings.Repeat("=", len(title)) + "\n" +
		title + "\n" +
		strings.Repeat("=", len(title)) + "\n" +
		"WorkerDirectory:      \t" + G.WorkerDirectory + "\n" +
		"LogDirectory:         \t" + G.LogDirectory + "\n" +
		"TemplateDirectory:    \t" + G.TemplateDirectory + "\n" +
		"TourneyDirectory:     \t" + G.TourneyDirectory + "\n" +
		"SaveDirectory:        \t" + G.SaveDirectory + "\n" +
		"BookDirectory:        \t" + G.BookDirectory + "\n" +
		"BookDirectory:        \t" + G.BuildDirectory + "\n" +
		"ServerPort:           \t" + strconv.Itoa(G.ServerPort) + "\n" +
		"WebPort:              \t" + strconv.Itoa(G.WebPort) + "\n" +
		"EngineFilePort:       \t" + strconv.Itoa(G.EngineFilePort) + "\n" +
		"MaxConnectionAttempts:\t" + strconv.Itoa(G.MaxConnectionAttempts) + "\n" +
		"\n" +
		"To modify these settings change them in the 'tourney.settings' file.\n"
}

func (G *GlobalSettings) Save(filename string) error {
	//check if the file exists:

	fmt.Print("Saving '" + filename + "'... ")
	var file *os.File
	var err error
	if _, er := os.Stat(filename); os.IsNotExist(er) {
		// file doesnt exist
	} else if er == nil {
		// file does exist
		os.Remove(filename)
	}

	file, err = os.Create(filename)
	defer file.Close()

	var encoded []byte
	encoded, err = json.MarshalIndent(G, "", "  ")
	//encoded, err = json.Marshal(G)
	if err != nil {
		return err
	}
	if _, err = file.Write(encoded); err != nil {
		return err
	}
	return nil
}

func (G *GlobalSettings) Load(filename string) error {

	var err error
	if _, err = os.Stat(filename); os.IsNotExist(err) {
		// file doesnt exist
		return err
	} else if err == nil {
		// file does exist
		file, err := os.Open(filename)
		defer file.Close()
		jsonParser := json.NewDecoder(file)
		if err = jsonParser.Decode(G); err != nil {
			return err
		}
	}
	fmt.Println("Program settings loaded from '" + filename + "'")
	return nil
}
