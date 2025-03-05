package application

import (
	"context"
	"github.com/sabure-dev/calc_go/internal/config"
	"github.com/sabure-dev/calc_go/internal/http/server"
	"github.com/sabure-dev/calc_go/internal/service/orchestrator"
	"os"
	"os/signal"
)

type Application struct {
	config *config.Config
}

func New() *Application {
	return &Application{
		config: config.ConfigFromEnv(),
	}
}

func (a *Application) Run(ctx context.Context) int {
	o := orchestrator.NewOrchestator()
	shutdownFunc := server.Run(a.config.Addr, o, a.config)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	ctx, cancel := context.WithCancel(context.Background())
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	<-c
	cancel()
	err := shutdownFunc(ctx)
	if err != nil {
		return 1
	}
	o.Shutdown()
	return 0
}
