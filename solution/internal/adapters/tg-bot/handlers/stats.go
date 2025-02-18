package handlers

import (
	"context"
	"fmt"
	"github.com/enescakir/emoji"
	"github.com/google/uuid"
	"github.com/nlypage/intele"
	"github.com/nlypage/intele/collector"
	"go.opentelemetry.io/otel"
	tele "gopkg.in/telebot.v3"
	"solution/cmd/app"
	v1 "solution/internal/adapters/controller/api/v1"
	"solution/internal/adapters/database/postgres"
	"solution/internal/adapters/logger"
	"solution/internal/domain/dto"
	"solution/internal/domain/service"
)

type StatsHandler struct {
	service      v1.StatsService
	inputManager *intele.InputManager
	bot          *tele.Bot
}

func NewStatsHandler(app *app.App, inputManager *intele.InputManager) *StatsHandler {
	return &StatsHandler{
		inputManager: inputManager,
		bot:          app.Telegram,
		service:      service.NewStatsService(postgres.NewStatsStorage(app.DB)),
	}
}

func (h *StatsHandler) StatsMenu(tgCtx tele.Context) error {
	buttons := []*tele.Btn{
		{Text: "Статистика по кампании", Unique: "stats_by_campaign"},
		{Text: "Статистика по рекламодателю", Unique: "stats_by_advertiser"},
		{Text: "Ежедневная статистика по кампании", Unique: "daily_stats_by_campaign"},
		{Text: "Ежедневная статистика по рекламодателю", Unique: "daily_stats_by_advertiser"},
	}

	err := tgCtx.Send(fmt.Sprintf("Выберите действие %v", emoji.BarChart), &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{*buttons[0].Inline()},
			{*buttons[1].Inline()},
			{*buttons[2].Inline()},
			{*buttons[3].Inline()},
		},
	})

	h.bot.Handle(buttons[0], h.StatsByCampaign)
	h.bot.Handle(buttons[1], h.StatsByAdvertiser)
	h.bot.Handle(buttons[2], h.DailyStatsByCampaign)
	h.bot.Handle(buttons[3], h.DailyStatsByAdvertiser)

	return err
}

func (h *StatsHandler) StatsByCampaign(tgCtx tele.Context) error {
	tracer := otel.Tracer("tg-stats-by-campaign-handler")
	ctx, span := tracer.Start(context.Background(), "Tg-StatsByCampaign")
	defer span.End()

	var statsDTO dto.GetStatsByCampaignIDDTO

	inputCollector := collector.New()

	_ = inputCollector.Send(tgCtx, "Введите ID кампании:")

	response, errGet := h.inputManager.Get(context.Background(), tgCtx.Sender().ID, 0)

	done := false

	for !done {
		switch {
		case response.Canceled:
			_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case errGet != nil || response.Message == nil:
			logger.Log.Errorf("failed to get input: %v", errGet)
			_ = inputCollector.Send(tgCtx, "Что-то пошло не так. Попробуем ещё раз")
		case uuid.Validate(response.Message.Text) != nil:
			_ = inputCollector.Send(tgCtx, "Некорректный ID Кампании, введите его ещё раз: ")
		case uuid.Validate(response.Message.Text) == nil:
			statsDTO.CampaignID = response.Message.Text
			_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
	}

	_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})

	stats, err := h.service.GetStatsByCampaignID(ctx, statsDTO)
	if err != nil {
		logger.Log.Errorf("failed to get stats: %v", err)
		return err
	}

	return tgCtx.Send(fmt.Sprintf(`
					Статистика по кампании:
					%v Количество показов: %v
					%v Количество кликов: %v
					%v Конверсия: %v
					%v Потрачено на просмотры: %v
					%v Потрачено на клики: %v
					%v Потрачено всего: %v`,
		emoji.BlackSmallSquare, stats.ImpressionsCount,
		emoji.BlackSmallSquare, stats.ClicksCount,
		emoji.BlackSmallSquare, stats.Conversion,
		emoji.BlackSmallSquare, stats.SpentImpressions,
		emoji.BlackSmallSquare, stats.SpentClicks,
		emoji.BlackSmallSquare, stats.SpentTotal))
}

func (h *StatsHandler) StatsByAdvertiser(tgCtx tele.Context) error {
	tracer := otel.Tracer("tg-stats-by-advertiser-handler")
	ctx, span := tracer.Start(context.Background(), "Tg-StatsByAdvertiser")
	defer span.End()

	var statsDTO dto.GetStatsByAdvertiserIDDTO

	inputCollector := collector.New()

	_ = inputCollector.Send(tgCtx, "Введите ID рекламодателя:")

	response, errGet := h.inputManager.Get(context.Background(), tgCtx.Sender().ID, 0)

	done := false

	for !done {
		switch {
		case response.Canceled:
			_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case errGet != nil || response.Message == nil:
			logger.Log.Errorf("failed to get input: %v", errGet)
			_ = inputCollector.Send(tgCtx, "Что-то пошло не так. Попробуем ещё раз")
		case uuid.Validate(response.Message.Text) != nil:
			_ = inputCollector.Send(tgCtx, "Некорректный ID Рекламодателя, введите его ещё раз: ")
		case uuid.Validate(response.Message.Text) == nil:
			statsDTO.AdvertiserID = response.Message.Text
			_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
	}

	_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})

	stats, err := h.service.GetStatsByAdvertiserID(ctx, statsDTO)
	if err != nil {
		logger.Log.Errorf("failed to get stats: %v", err)
		return err
	}

	return tgCtx.Send(fmt.Sprintf(`
					Статистика по всем кампаниям рекламодателя:
					%v Количество показов: %v
					%v Количество кликов: %v
					%v Конверсия: %v
					%v Потрачено на просмотры: %v
					%v Потрачено на клики: %v
					%v Потрачено всего: %v`,
		emoji.BlackSmallSquare, stats.ImpressionsCount,
		emoji.BlackSmallSquare, stats.ClicksCount,
		emoji.BlackSmallSquare, stats.Conversion,
		emoji.BlackSmallSquare, stats.SpentImpressions,
		emoji.BlackSmallSquare, stats.SpentClicks,
		emoji.BlackSmallSquare, stats.SpentTotal))
}

func (h *StatsHandler) DailyStatsByCampaign(tgCtx tele.Context) error {
	tracer := otel.Tracer("tg-stats-by-campaign-handler")
	ctx, span := tracer.Start(context.Background(), "Tg-StatsByCampaign")
	defer span.End()

	var statsDTO dto.GetStatsByCampaignIDDTO

	inputCollector := collector.New()

	_ = inputCollector.Send(tgCtx, "Введите ID кампании:")

	response, errGet := h.inputManager.Get(context.Background(), tgCtx.Sender().ID, 0)

	done := false

	for !done {
		switch {
		case response.Canceled:
			_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case errGet != nil || response.Message == nil:
			logger.Log.Errorf("failed to get input: %v", errGet)
			_ = inputCollector.Send(tgCtx, "Что-то пошло не так. Попробуем ещё раз")
		case uuid.Validate(response.Message.Text) != nil:
			_ = inputCollector.Send(tgCtx, "Некорректный ID Кампании, введите его ещё раз: ")
		case uuid.Validate(response.Message.Text) == nil:
			statsDTO.CampaignID = response.Message.Text
			_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
	}

	_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})

	stats, err := h.service.GetDailyStatsByCampaignID(ctx, statsDTO)
	if err != nil {
		logger.Log.Errorf("failed to get stats: %v", err)
		return err
	}

	for _, stat := range stats {
		_ = tgCtx.Send(fmt.Sprintf(`
					Статистика по кампании за день %v:
					%v Количество показов: %v
					%v Количество кликов: %v
					%v Конверсия: %v
					%v Потрачено на просмотры: %v
					%v Потрачено на клики: %v
					%v Потрачено всего: %v`, stat.Day,
			emoji.BlackSmallSquare, stat.ImpressionsCount,
			emoji.BlackSmallSquare, stat.ClicksCount,
			emoji.BlackSmallSquare, stat.Conversion,
			emoji.BlackSmallSquare, stat.SpentImpressions,
			emoji.BlackSmallSquare, stat.SpentClicks,
			emoji.BlackSmallSquare, stat.SpentTotal))
	}

	return nil
}

func (h *StatsHandler) DailyStatsByAdvertiser(tgCtx tele.Context) error {
	tracer := otel.Tracer("tg-stats-by-advertiser-handler")
	ctx, span := tracer.Start(context.Background(), "Tg-StatsByAdvertiser")
	defer span.End()

	var statsDTO dto.GetStatsByAdvertiserIDDTO

	inputCollector := collector.New()

	_ = inputCollector.Send(tgCtx, "Введите ID Рекламодателя:")

	response, errGet := h.inputManager.Get(context.Background(), tgCtx.Sender().ID, 0)

	done := false

	for !done {
		switch {
		case response.Canceled:
			_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
			return nil
		case errGet != nil || response.Message == nil:
			logger.Log.Errorf("failed to get input: %v", errGet)
			_ = inputCollector.Send(tgCtx, "Что-то пошло не так. Попробуем ещё раз")
		case uuid.Validate(response.Message.Text) != nil:
			_ = inputCollector.Send(tgCtx, "Некорректный ID Рекламодателя, введите его ещё раз: ")
		case uuid.Validate(response.Message.Text) == nil:
			statsDTO.AdvertiserID = response.Message.Text
			_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
	}

	_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})

	stats, err := h.service.GetDailyStatsByAdvertiserID(ctx, statsDTO)
	if err != nil {
		logger.Log.Errorf("failed to get stats: %v", err)
		return err
	}

	for _, stat := range stats {
		_ = tgCtx.Send(fmt.Sprintf(`
					Статистика по всем кампаниям рекламодателя за день %v:
					%v Количество показов: %v
					%v Количество кликов: %v
					%v Конверсия: %v
					%v Потрачено на просмотры: %v
					%v Потрачено на клики: %v
					%v Потрачено всего: %v`, stat.Day,
			emoji.BlackSmallSquare, stat.ImpressionsCount,
			emoji.BlackSmallSquare, stat.ClicksCount,
			emoji.BlackSmallSquare, stat.Conversion,
			emoji.BlackSmallSquare, stat.SpentImpressions,
			emoji.BlackSmallSquare, stat.SpentClicks,
			emoji.BlackSmallSquare, stat.SpentTotal))
	}

	return nil
}
