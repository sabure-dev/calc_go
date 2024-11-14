package application

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
}

type Application struct {
	config *Config
	server *http.Server
}

func New() *Application {
	return &Application{
		config: loadConfig(),
	}
}

func loadConfig() *Config {
	config := &Config{}

	log.Println("Загрузка файла конфигурации...")
	err := godotenv.Load()
	if err != nil {
		log.Println("Ошибка загрузки файла конфигурации:", err)
	} else {
		log.Println("Файл конфигурации загружен успешно")
	}

	config.Port = os.Getenv("PORT")

	if config.Port == "" {
		log.Println("PORT не установлен, используется порт по умолчанию 8080")
		config.Port = "8080"
	}

	return config
}

func (a *Application) Run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/calculate", a.calculateHandler)

	handler := loggingMiddleware(mux)

	a.server = &http.Server{
		Addr:    ":" + a.config.Port,
		Handler: handler,
	}

	log.Printf("Запуска сервера на порту %s\n", a.config.Port)
	return a.server.ListenAndServe()
}
