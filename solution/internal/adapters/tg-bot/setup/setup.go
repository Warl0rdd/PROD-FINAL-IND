package setup

import (
	"github.com/nlypage/intele"
	"go.uber.org/zap"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"solution/cmd/app"
	"solution/internal/adapters/logger"
	"solution/internal/adapters/tg-bot/handlers"
)

func Setup(app *app.App) {
	bot := app.Telegram

	bot.Use(middleware.Recover(func(err error, ctx tele.Context) {
		logger.Log.Panicf("Telegram bot panic! %v", err)
	}))

	bot.Use(middleware.Logger(zap.NewStdLog(logger.Log.Desugar())))

	inputManager := intele.NewInputManager(intele.InputOptions{})

	bot.Handle(tele.OnText, inputManager.MessageHandler())
	bot.Handle(tele.OnCallback, inputManager.CallbackHandler())

	startHandler := handlers.NewStartHandler(inputManager,
		bot,
		handlers.NewCampaignHandler(app, inputManager),
		handlers.NewStatsHandler(app, inputManager))
	startHandler.Setup(bot)

	go bot.Start()
}
