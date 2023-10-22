package main

import (
	"context"
	"refactoring/internal/app"
	"refactoring/internal/app/startup"
	"refactoring/internal/http"
	"refactoring/internal/service"
)

func main() {

	// Файл с конфигурацией проекта
	configPath := "config/config.yaml"

	// Создаём логгер
	logger := startup.NewLogger()

	// Парсим файл конфигурации
	config, err := startup.NewConfig(configPath)
	if err != nil {
		logger.Fatalf("failed to Config: %v", err)
	}

	// Клиент для реализации бизнес-логики
	client := service.NewService(logger, config.Store.FilePath)
	err = client.Start()
	if err != nil {
		logger.Fatalf("failed Start Service: %v", err)
	}

	// Создаём экземпляр http сервера
	httpRouter := http.NewHttpRouter(config.Http, client, logger)

	// Запускаем http сервер
	app.NewApp(logger, httpRouter).Run(context.Background())
}
