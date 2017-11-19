// Package worker plays games. Games are fetched from the server.
package worker

import (
	"fmt"
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/worker/client"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	id     models.Id
	client *client.ApiClient
}

func New(apiURL string) *Worker {
	return &Worker{
		id:     models.NewId(),
		client: client.New(apiURL),
	}
}

func (w *Worker) Start() {
	fmt.Println("Starting worker.")
	g, err := w.client.NextGame()
	if err != nil {
		panic(err)
	}
	log.Debug("Recieved game: ", g)
	w.claim(g)
	w.getEngines(g)
	w.play(g)
}
