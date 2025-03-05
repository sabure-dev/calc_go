package core

import (
	"sync"
	"time"
)

type Task struct {
	Id            int     `json:"id"`
	Operation     string  `json:"operation"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Result        float64
	OperationTime time.Duration `json:"operation_time"`
}

type TaskWrapper struct {
	Task Task `json:"task"`
}

func Worker(tasks <-chan Task, results chan<- Task, wg *sync.WaitGroup) {
	defer wg.Done()
	for task := range tasks {
		var res float64
		switch task.Operation {
		case "+":
			res = task.Arg1 + task.Arg2
		case "-":
			res = task.Arg1 - task.Arg2
		case "*":
			res = task.Arg1 * task.Arg2
		case "/":
			res = task.Arg1 / task.Arg2
		}
		task.Result = res
		time.Sleep(task.OperationTime)
		results <- task
	}
}
