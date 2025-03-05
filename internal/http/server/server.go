package server

import (
	"context"
	"github.com/sabure-dev/calc_go/internal/config"
	"github.com/sabure-dev/calc_go/internal/http/server/handler"
	"github.com/sabure-dev/calc_go/internal/service/orchestrator"
	"github.com/sabure-dev/calc_go/pkg/calc/safeStructs"
	"go.uber.org/zap"
	"net/http"
)

func newHandler(o *orchestrator.Orchestator, config *config.Config) http.Handler {
	muxHandler := http.NewServeMux()
	Map := safeStructs.NewSafeMap()
	Id := safeStructs.NewSafeId()
	muxHandler.HandleFunc("/api/v1/calculate", func(w http.ResponseWriter, r *http.Request) {
		handler.CalcHandler(w, r, o, Map, Id)
	})
	muxHandler.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		handler.GiveTask(w, r, o, config.Delay, Map)
	})
	muxHandler.HandleFunc("/api/v1/expressions/", func(w http.ResponseWriter, r *http.Request) {
		handler.GetExpression(w, r, Map)
	})
	muxHandler.HandleFunc("/api/v1/expressions", func(w http.ResponseWriter, r *http.Request) {
		handler.GetAllExpressions(w, r, Map)
	})
	return handler.Decorate(muxHandler)
}

func Run(addr string, o *orchestrator.Orchestator, configVar *config.Config) func(ctx context.Context) error {
	Handler := newHandler(o, configVar)
	server := &http.Server{Addr: ":" + addr, Handler: Handler}
	ch := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		if err != nil {
			ch <- err
		}
	}()
	select {
	case err := <-ch:
		if err != nil {
			zap.String("Err", err.Error())
		}
		return server.Shutdown
	}
}
