package setup

import (
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/logger"
	"github.com/spf13/viper"
	"solution/cmd/app"
	v1 "solution/internal/adapters/controller/api/v1"
)

func Setup(app *app.App) {
	app.Fiber.Use(cors.New(cors.ConfigDefault))

	if viper.GetBool("settings.debug") {
		app.Fiber.Use(logger.New(logger.Config{TimeZone: viper.GetString("settings.timezone")}))
	}

	router := app.Fiber.Group("")

	clientHandler := v1.NewClientHandler(app)
	clientHandler.Setup(router)

	advertiserHandler := v1.NewAdvertiserHandler(app)
	advertiserHandler.Setup(router)

	mlScoreHandler := v1.NewMlScoreHandler(app)
	mlScoreHandler.Setup(router)

	campaignHandler := v1.NewCampaignHandler(app)
	campaignHandler.Setup(router)

	dayHandler := v1.NewDayHandler(app)
	dayHandler.Setup(router)

	adsHandler := v1.NewAdsHandler(app)
	adsHandler.Setup(router)

	statsHandler := v1.NewStatsHandler(app)
	statsHandler.Setup(router)

	imageHandler := v1.NewImageHandler(app)
	imageHandler.Setup(router)

	moderationHandler := v1.NewModerationHandler(app)
	moderationHandler.Setup(router)

	llmHandler := v1.NewLLMHandler(app)
	llmHandler.Setup(router)
}
