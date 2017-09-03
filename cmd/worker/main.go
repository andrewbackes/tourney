package main

import (
	"fmt"
	"github.com/andrewbackes/tourney/worker"
	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"
	"sync"
)

const apiURL = "http://localhost:9090/api/v2"
const workerCount = 1

func main() {
	log.SetFormatter(new(prefixed.TextFormatter))
	log.SetLevel(log.DebugLevel)

	fmt.Println("Worker.")
	workers := make([]*worker.Worker, 0)
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		w := worker.New(apiURL)
		workers = append(workers, w)
		wg.Add(1)
		go func() {
			w.Start()
			wg.Done()
		}()
	}
	wg.Wait()
}
