package main

import (
	"solution/cmd/app"
	"solution/internal/adapters/config"
	"solution/internal/adapters/controller/api/setup"
	_ "solution/internal/adapters/controller/api/v1"
	botSetup "solution/internal/adapters/tg-bot/setup"
)

func main() {
	appConfig := config.Configure(false)
	mainApp := app.New(appConfig)

	defer mainApp.DB.Close()

	setup.Setup(mainApp)

	botSetup.Setup(mainApp)

	mainApp.Start()
}
