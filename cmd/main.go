package main

import (
	"log"

	"github.com/sabure-dev/calc_go/internal/application"
)

func main() {
	app := application.New()
	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
