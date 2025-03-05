package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/sabure-dev/calc_go/internal/agent/core"
	"net/http"
	"os"
	"strconv"
	"time"
)

type Transport struct {
	TaskChan   chan core.Task
	ResultChan chan core.Task
}

func NewAgent(workerCount int) *Transport {
	return &Transport{
		TaskChan:   make(chan core.Task, workerCount),
		ResultChan: make(chan core.Task, workerCount),
	}
}

func (t *Transport) Start(serverURL string, ctx context.Context) {

	pingDelay := getPingDelay()

	go t.fetchTasks(serverURL, pingDelay, ctx)

	go t.sendResults(serverURL, ctx)
}

func (t *Transport) fetchTasks(serverURL string, pingDelay time.Duration, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			time.Sleep(pingDelay)

			resp, err := http.Get(serverURL)
			if err != nil {
				continue
			}
			defer resp.Body.Close()

			if resp.StatusCode == http.StatusNotFound || resp.StatusCode == http.StatusInternalServerError {
				continue
			}

			var taskWrapper core.TaskWrapper
			if err := json.NewDecoder(resp.Body).Decode(&taskWrapper); err != nil {
				continue
			}

			t.TaskChan <- taskWrapper.Task
		}
	}
}

func (t *Transport) sendResults(serverURL string, ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case resultTask := <-t.ResultChan:
			taskWrapper := core.TaskWrapper{Task: resultTask}

			jsonData, err := json.Marshal(taskWrapper)
			if err != nil {
				continue
			}

			_, err = http.Post(serverURL, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
			}
		}
	}
}

func getPingDelay() time.Duration {
	pingStr := os.Getenv("PING")
	if pingStr == "" {
		pingStr = "1000"
	}
	pingDelay, _ := strconv.Atoi(pingStr)

	return time.Duration(pingDelay) * time.Millisecond
}

func (t *Transport) Shutdown() {
	close(t.TaskChan)
	close(t.ResultChan)
}
