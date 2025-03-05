package main

import (
	"context"
	"github.com/sabure-dev/calc_go/internal/agent/core"
	"github.com/sabure-dev/calc_go/internal/agent/transport"
	"log"
	"os"
	"os/signal"
	"strconv"
	"sync"
)

func main() {
	goroutinesEnv := os.Getenv("COMPUTING_POWER")
	if goroutinesEnv == "" {
		goroutinesEnv = "2"
	}

	numWorkers, err := strconv.Atoi(goroutinesEnv)
	if err != nil {
		log.Fatalf("Error: %v", err)
	}

	calcAgent := transport.NewAgent(numWorkers)

	serverURL := os.Getenv("URL")
	if serverURL == "" {
		serverURL = "http://127.0.0.1:8080"
	}
	taskEndpoint := serverURL + "/internal/task"

	workerGroup := &sync.WaitGroup{}
	for i := 0; i < numWorkers; i++ {
		workerGroup.Add(1)
		go core.Worker(calcAgent.TaskChan, calcAgent.ResultChan, workerGroup)
	}

	ctx, cancel := context.WithCancel(context.Background())

	calcAgent.Start(taskEndpoint, ctx)

	interruptChan := make(chan os.Signal, 1)
	signal.Notify(interruptChan, os.Interrupt)
	<-interruptChan

	cancel()
	calcAgent.Shutdown()
	os.Exit(0)
}
