package main

import (
	"context"
	"log"
	"refactoring/internal/app"
	"refactoring/internal/app/startup"
	"refactoring/internal/http"
	"refactoring/internal/service"
)

const loggerName = "UserApi"

func main() {

	//configPath := os.Getenv("CONFIG_PATH")

	// Файл с конфигурацией проекта
	configPath := "config/config.local.yaml"

	// Парсим файл конфигурации
	config, err := startup.NewConfig(configPath)
	if err != nil {
		log.Fatalf("failed to Config %v", err)
	}

	// Создаём логгер
	logger := startup.NewLogger(loggerName)

	// Клиент для реализации бизнес-логики
	client := service.NewService(logger, config.Store.FilePath)
	err = client.Start()
	if err != nil {
		log.Fatalf("failed Start Service: %v", err)
	}

	// Создаём экземпляр http сервера
	httpRouter := http.NewHttpRouter(config.Http, client, logger)

	// Запускаем http сервер
	app.NewApp(logger, httpRouter).Run(context.Background())
}
