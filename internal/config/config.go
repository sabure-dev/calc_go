package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

type Delay struct {
	Plus     time.Duration
	Minus    time.Duration
	Multiple time.Duration
	Divide   time.Duration
}

type Config struct {
	Addr  string
	Delay Delay
	Mode  Mode
}

type Mode struct {
	Console   string
	File      string
	CleanFile string
}

func configDelay(config *Config) {
	plus := os.Getenv("TIME_ADDITION_MS")
	if plus == "" {
		config.Delay.Plus = time.Second
	} else {
		plus, err := strconv.Atoi(plus)
		if err != nil {
			log.Fatalf("Ошибка преобразования TIME_ADDITION_MS: %s", err)
		}
		config.Delay.Plus = time.Duration(plus) * time.Millisecond
	}
	minus := os.Getenv("TIME_SUBTRACTION_MS")
	if minus == "" {
		config.Delay.Minus = time.Second
	} else {
		minus, err := strconv.Atoi(minus)
		if err != nil {
			log.Fatalf("Ошибка преобразования TIME_SUBTRACTION_MS: %s", err)
		}
		config.Delay.Minus = time.Duration(minus) * time.Millisecond
	}
	multiple := os.Getenv("TIME_MULTIPLICATIONS_MS")
	if multiple == "" {
		config.Delay.Multiple = time.Second
	} else {
		multiple, err := strconv.Atoi(multiple)
		if err != nil {
			log.Fatalf("Ошибка преобразования TIME_MULTIPLICATION_MS: %s", err)
		}
		config.Delay.Multiple = time.Duration(multiple) * time.Millisecond
	}
	divide := os.Getenv("TIME_DIVISIONS_MS")
	if divide == "" {
		config.Delay.Divide = time.Second
	} else {
		divide, err := strconv.Atoi(divide)
		if err != nil {
			log.Fatalf("Ошибка преобразования TIME_DIVISION_MS: %s", err)
		}
		config.Delay.Divide = time.Duration(divide) * time.Millisecond
	}
}

func ConfigFromEnv() *Config {
	config := new(Config)
	config.Addr = "8080"

	config.Mode.Console = os.Getenv("MODE_CONSOLE")
	if config.Mode.Console == "" {
		config.Mode.Console = "Dev"
	}
	config.Mode.File = os.Getenv("MODE_FILE")
	if config.Mode.File == "" {
		config.Mode.File = "Prod"
	}
	config.Mode.CleanFile = os.Getenv("CLEAN_FILE")
	if config.Mode.CleanFile == "" {
		config.Mode.CleanFile = "False"
	}
	configDelay(config)
	return config
}
