package handler

import (
	"encoding/json"
	"fmt"
	"github.com/sabure-dev/calc_go/internal/agent/core"
	"github.com/sabure-dev/calc_go/internal/config"
	"github.com/sabure-dev/calc_go/internal/service/orchestrator"
	"github.com/sabure-dev/calc_go/pkg/calc"
	"github.com/sabure-dev/calc_go/pkg/calc/safeStructs"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Decorator func(http.Handler) http.Handler

type Request struct {
	Expression string `json:"expression"`
}

type ResponseWr struct {
	Expression safeStructs.Expressions `json:"expression"`
}

type ExprWr struct {
	Expressions []safeStructs.Expressions `json:"expressions"`
}

type ResoponseId struct {
	ID int `json:"id"`
}

type ResultBad struct {
	Err string `json:"error"`
}

func CalcHandler(w http.ResponseWriter, r *http.Request, o *orchestrator.Orchestator, Map *safeStructs.SafeMap, Id *safeStructs.SafeId) {
	request := new(Request)
	err := json.NewDecoder(r.Body).Decode(&request)
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	if err != nil && err != io.EOF {
		w.WriteHeader(422)
		errj := calc.ErrInvalidJson
		res := ResultBad{Err: errj.Error()}
		jsonBytes, _ := json.Marshal(res)
		fmt.Fprint(w, string(jsonBytes))
		time.Sleep(1)
		return
	} else if err == io.EOF {
		w.WriteHeader(422)
		errj := calc.ErrEmptyExpression
		res := ResultBad{Err: errj.Error()}
		jsonBytes, _ := json.Marshal(res)
		fmt.Fprint(w, string(jsonBytes))
		time.Sleep(1)
		return
	} else {
		w.WriteHeader(201)
	}
	id := Id.Get()
	Map.Set(id, safeStructs.Expressions{Id: id, Status: "Подсчёт"})
	resp := ResoponseId{ID: id}
	jsonBytes, err := json.Marshal(resp)
	if err != nil {
		w.WriteHeader(500)
	} else {
		fmt.Fprint(w, string(jsonBytes))
	}
	go o.Calculate(request.Expression, id, Map)
	time.Sleep(1)
}

func GiveTask(w http.ResponseWriter, r *http.Request, o *orchestrator.Orchestator, delay config.Delay, safeMap *safeStructs.SafeMap) {
	if r.Method == "GET" {
		if len(o.Out) == 0 {
			w.WriteHeader(404)
		} else {
			task := <-o.Out
			switch task.Operation {
			case "+":
				task.OperationTime = delay.Plus
			case "-":
				task.OperationTime = delay.Minus
			case "*":
				task.OperationTime = delay.Multiple
			case "/":
				task.OperationTime = delay.Divide
			}
			taskWr := core.TaskWrapper{Task: task}
			jsonBytes, err := json.Marshal(taskWr)
			if err != nil {
				w.WriteHeader(500)
			}
			fmt.Fprint(w, string(jsonBytes))
		}
	} else {
		taskRes := new(core.TaskWrapper)
		err := json.NewDecoder(r.Body).Decode(taskRes)
		defer r.Body.Close()
		if err != nil {
			w.WriteHeader(500)
		} else {
			task := taskRes.Task
			if !safeMap.In(task.Id) {
				w.WriteHeader(404)
			} else if safeMap.Get(task.Id).Result != "" {
				w.WriteHeader(422)
			} else {
				o.In <- task.Result
				w.WriteHeader(200)
			}
		}
	}
}

func GetExpression(w http.ResponseWriter, r *http.Request, safeMap *safeStructs.SafeMap) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	url := r.URL.Path
	idStr := strings.TrimPrefix(url, "/api/v1/expressions/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(500)
		return
	}
	res := safeMap.Get(id)
	if res.Id == 0 {
		w.WriteHeader(404)
	} else {
		w.WriteHeader(200)
		resWr := ResponseWr{Expression: res}
		jsonBytes, err := json.Marshal(resWr)
		if err != nil {
			w.WriteHeader(500)
			return
		}
		fmt.Fprint(w, string(jsonBytes))
	}
}

func GetAllExpressions(w http.ResponseWriter, r *http.Request, safeMap *safeStructs.SafeMap) {
	if r.Method != http.MethodGet {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	expressions := safeMap.GetAll()
	res := ExprWr{Expressions: expressions}
	jsonBytes, err := json.Marshal(res)
	if err != nil {
		w.WriteHeader(500)
	} else {
		fmt.Fprint(w, string(jsonBytes))
	}
}

func Decorate(next http.Handler, ds ...Decorator) http.Handler {
	res := next
	for d := len(ds) - 1; d >= 0; d-- {
		res = ds[d](res)
	}
	return res
}
