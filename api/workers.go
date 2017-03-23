package api

import (
	"github.com/andrewbackes/tourney/model/structures"
	"net/http"
)

func (c *controller) postWorker(w http.ResponseWriter, req *http.Request) {
	var worker structures.Worker
	readJSON(req.Body, &worker)
	if worker.Id == "" {
		worker = *structures.NewWorker()
	}
	c.model.AddWorker(&worker)
	writeJSON(&worker, w)
}
