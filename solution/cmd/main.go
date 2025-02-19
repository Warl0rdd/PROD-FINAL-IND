package main

import (
	"solution/cmd/app"
	"solution/internal/adapters/config"
	"solution/internal/adapters/controller/api/setup"
	_ "solution/internal/adapters/controller/api/v1"
	botSetup "solution/internal/adapters/tg-bot/setup"
)

// @title           Опциональные требования
// @version         1.0
// @description     Описание работы всех API эндпоинтов из опциональных требований

// @host      localhost:8080
// @BasePath  /
func main() {
	appConfig := config.Configure()
	mainApp := app.New(appConfig)

	defer mainApp.DB.Close()

	setup.Setup(mainApp)

	botSetup.Setup(mainApp)

	mainApp.Start()
}
