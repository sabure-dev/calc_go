package main

import (
	"context"
	"github.com/sabure-dev/calc_go/internal/application"
	"os"
)

func main() {
	app := application.New()
	ctx := context.Background()

	os.Exit(app.Run(ctx))
}
