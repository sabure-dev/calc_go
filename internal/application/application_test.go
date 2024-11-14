package application

import (
	"net/http"
	"os"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	os.Setenv("PORT", "9090")

	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}

	if app.config == nil {
		t.Fatal("config is nil")
	}

	if app.config.Port != "9090" {
		t.Errorf("expected port 9090, got %s", app.config.Port)
	}
}

func TestLoadConfig(t *testing.T) {
	originalPort := os.Getenv("PORT")
	defer os.Setenv("PORT", originalPort)

	tests := []struct {
		name     string
		envPort  string
		expected string
	}{
		{
			name:     "with PORT env",
			envPort:  "3000",
			expected: "3000",
		},
		{
			name:     "without PORT env",
			envPort:  "",
			expected: "8080",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("PORT", tt.envPort)
			config := loadConfig()

			if config.Port != tt.expected {
				t.Errorf("expected port %s, got %s", tt.expected, config.Port)
			}
		})
	}
}

func TestRun(t *testing.T) {
	os.Setenv("PORT", "8082")

	app := New()

	errChan := make(chan error, 1)

	go func() {
		errChan <- app.Run()
	}()

	time.Sleep(100 * time.Millisecond)

	select {
	case err := <-errChan:
		if err != nil && err != http.ErrServerClosed {
			t.Errorf("unexpected error: %v", err)
		}
	default:
	}
}
