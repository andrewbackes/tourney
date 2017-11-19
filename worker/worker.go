// Package worker plays games. Games are fetched from the server.
package worker

import (
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
	log.Info("Starting worker.")
	g, err := w.client.NextPendingGame()
	if err != nil {
		panic(err)
	}
	log.Debug("Recieved game: ", g)
	w.claim(g)
	w.getEngines(g)
	w.play(g)
}
