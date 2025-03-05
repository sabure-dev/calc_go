package orchestrator

import (
	"errors"
	"github.com/sabure-dev/calc_go/internal/agent/core"
	Calc "github.com/sabure-dev/calc_go/internal/service/orchestrator/core"
	"github.com/sabure-dev/calc_go/pkg/calc"
	"github.com/sabure-dev/calc_go/pkg/calc/safeStructs"
	"strconv"
)

type Orchestator struct {
	Out      chan core.Task
	In       chan float64
	ErrorsCh chan error
	Ready    chan int
}

func NewOrchestator() *Orchestator {
	return &Orchestator{
		Out:      make(chan core.Task, 128),
		In:       make(chan float64, 128),
		ErrorsCh: make(chan error, 128),
		Ready:    make(chan int, 1),
	}
}

func (o *Orchestator) TakeExpression(expression string, id int) {
	Calc.Calc(expression, o.Out, o.In, o.ErrorsCh, o.Ready, id)
}

func (o *Orchestator) Calculate(expression string, id int, Map *safeStructs.SafeMap) {
	o.TakeExpression(expression, id)
	<-o.Ready
	var result float64
	var err error
	if len(o.ErrorsCh) > 0 {
		err = <-o.ErrorsCh
	} else {
		result = <-o.In
	}
	if err != nil {
		if _, ok := calc.ErrorMap[err]; ok {
			Map.Set(id, safeStructs.Expressions{Id: id, Status: "Выполнено", Result: err.Error()})
		} else {
			errJ := errors.New("Что-то пошло не так")
			Map.Set(id, safeStructs.Expressions{Id: id, Status: "Выполнено", Result: errJ.Error()})
		}
	} else {
		resStr := strconv.FormatFloat(result, 'f', 2, 64)
		Map.Set(id, safeStructs.Expressions{Id: id, Status: "Выполнено", Result: resStr})
	}
}

func (o *Orchestator) Shutdown() {
	close(o.Out)
	close(o.In)
	close(o.ErrorsCh)
	close(o.Ready)
}
