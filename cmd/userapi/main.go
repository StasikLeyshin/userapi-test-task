package main

import (
	"log"
	"refactoring/internal/app/startup"
	"refactoring/internal/service"
)

const loggerName = "UserApi"

func main() {

	//configPath := os.Getenv("CONFIG_PATH")

	configPath := "config/config.local.yaml"

	config, err := startup.NewConfig(configPath)
	if err != nil {
		log.Fatalf("failed to Config %v", err)
	}

	logger := startup.NewLogger(loggerName)

	// Клиент для реализации бизнес-логики
	client := service.NewService(logger, config.Store.FilePath)
	err = client.Start()
	if err != nil {
		log.Fatalf("failed Start Service: %v", err)
	}

	//// Запуск http сервера
	//httpRouter := http.NewHttpRouter(config.Http, client)
	//log.Printf("start http-server")
	//err = httpRouter.Start()
	//if err != nil {
	//	log.Fatalf("failed to listen the port: %d %v", config.Http.Port, err)
	//}
}
