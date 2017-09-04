// Package worker plays games. Games are fetched from the server.
package worker

import (
	"fmt"
	"github.com/andrewbackes/tourney/data/models"
	"github.com/andrewbackes/tourney/worker/services/master"
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	id     models.Id
	master *master.MasterService
}

func New(apiURL string) *Worker {
	return &Worker{
		id:     models.NewId(),
		master: master.New(apiURL),
	}
}

func (w *Worker) Start() {
	fmt.Println("Starting worker.")
	g, err := w.master.NextGame()
	if err != nil {
		panic(err)
	}
	log.Debug("Recieved game: ", g)
	w.claim(g)
	w.play(g)
}
