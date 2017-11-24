package main

import (
	"fmt"
	"github.com/andrewbackes/tourney/worker"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"os"
	"sync"
)

const workerCount = 1

func main() {
	log.SetFormatter(new(prefixed.TextFormatter))
	log.SetLevel(log.DebugLevel)

	fmt.Println("Worker.")
	workers := make([]*worker.Worker, 0)
	var wg sync.WaitGroup
	APIURL := getAPIURL()
	log.Info("Using API URL: ", APIURL)
	for i := 0; i < workerCount; i++ {
		w := worker.New(APIURL)
		workers = append(workers, w)
		wg.Add(1)
		go func() {
			w.Start()
			wg.Done()
		}()
	}
	wg.Wait()
}

func getAPIURL() string {
	if os.Getenv("API_URL") != "" {
		return os.Getenv("API_URL")
	}
	return "http://api.tourney.aback.es:9090/api/v2"
}
