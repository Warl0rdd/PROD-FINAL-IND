package main

import (
	"solution/cmd/app"
	"solution/internal/adapters/config"
	"solution/internal/adapters/controller/api/setup"
	botSetup "solution/internal/adapters/tg-bot/setup"
)

func main() {
	appConfig := config.Configure()
	mainApp := app.New(appConfig)

	defer mainApp.DB.Close()

	setup.Setup(mainApp)

	botSetup.Setup(mainApp)

	mainApp.Start()
}
