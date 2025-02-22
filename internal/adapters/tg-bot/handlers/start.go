package handlers

import (
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/nlypage/intele"
	tele "gopkg.in/telebot.v3"
	"solution/internal/adapters/logger"
)

type StartHandler struct {
	inputManager *intele.InputManager
	bot          *tele.Bot

	campaignHandler *CampaignHandler
	statsHandler    *StatsHandler
}

func NewStartHandler(inputManager *intele.InputManager, bot *tele.Bot, campaignHandler *CampaignHandler, statsHandler *StatsHandler) *StartHandler {
	return &StartHandler{
		inputManager:    inputManager,
		bot:             bot,
		campaignHandler: campaignHandler,
		statsHandler:    statsHandler,
	}
}

var startMessage = fmt.Sprintf(`
Привет! %v

Я - бот для управления рекламными кампаниями %v

Выберите действие:
`, emoji.WavingHand, emoji.Laptop)

func (h *StartHandler) Start(tgCtx tele.Context) error {
	buttons := []*tele.Btn{
		{Text: "Кампании", Unique: "campaigns"},
		{Text: "Статистика", Unique: "stats"},
	}

	err := tgCtx.Send(startMessage, &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{*buttons[0].Inline()},
			{*buttons[1].Inline()},
		},
	})

	logger.Log.Debugf("Sent start message to user: %v", tgCtx.Sender().ID)

	if err != nil {
		logger.Log.Errorf("failed to send message: %v", err)
	}

	h.bot.Handle(buttons[0], h.campaignHandler.CampaignMenu)

	h.bot.Handle(buttons[1], h.statsHandler.StatsMenu)

	return err
}

func (h *StartHandler) Setup(bot *tele.Bot) {
	bot.Handle("/start", h.Start)
}
