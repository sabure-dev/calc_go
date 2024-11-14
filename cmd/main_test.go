package main

import (
	"os"
	"testing"

	"github.com/sabure-dev/calc_go/internal/application"
)

func TestMain(t *testing.T) {
	os.Setenv("PORT", "8081")

	errChan := make(chan error, 1)

	go func() {
		app := application.New()
		errChan <- app.Run()
	}()

	select {
	case err := <-errChan:
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	default:
	}
}
