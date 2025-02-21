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
	"solution/internal/adapters/database/redis"
	"solution/internal/adapters/logger"
	"solution/internal/domain/dto"
	"solution/internal/domain/service"
	"solution/internal/domain/utils/parsing"
	"solution/internal/domain/utils/pointers"
	"strconv"
)

type CampaignHandler struct {
	service      v1.CampaignService
	inputManager *intele.InputManager
	bot          *tele.Bot
}

func NewCampaignHandler(app *app.App, inputManager *intele.InputManager) *CampaignHandler {
	return &CampaignHandler{
		service:      service.NewCampaignService(postgres.NewCampaignStorage(app.DB), redis.NewDayStorage(app.Redis)),
		inputManager: inputManager,
	}
}

func (h *CampaignHandler) CampaignMenu(tgCtx tele.Context) error {
	buttons := []*tele.Btn{
		{Text: "Создать кампанию", Unique: "create_campaign"},
		{Text: "Просмотреть рекламные кампании", Unique: "get_all"},
		{Text: "Обновить кампанию", Unique: "update_campaign"},
		{Text: "Удалить кампанию", Unique: "delete_campaign"},
	}

	err := tgCtx.Send(fmt.Sprintf("Выберите действие %v", emoji.Laptop), &tele.ReplyMarkup{
		InlineKeyboard: [][]tele.InlineButton{
			{*buttons[0].Inline()},
			{*buttons[1].Inline()},
			{*buttons[2].Inline()},
			{*buttons[3].Inline()},
		},
	})

	h.bot.Handle(buttons[0], func(tgCtx tele.Context) error {
		return h.CreateCampaign(context.Background(), tgCtx)
	})

	h.bot.Handle(buttons[1], func(tgCtx tele.Context) error {
		return h.GetAll(context.Background(), tgCtx)
	})

	h.bot.Handle(buttons[2], func(tgCtx tele.Context) error {
		return h.UpdateCampaign(context.Background(), tgCtx)
	})

	h.bot.Handle(buttons[3], func(tgCtx tele.Context) error {
		return h.DeleteCampaign(context.Background(), tgCtx)
	})

	return err
}

func (h *CampaignHandler) CreateCampaign(traceCtx context.Context, tgCtx tele.Context) error {
	tracer := otel.Tracer("campaign-handler")
	traceCtx, span := tracer.Start(traceCtx, "Tg-CreateCampaign")
	defer span.End()

	var createCampaignDTO dto.CreateCampaignDTO

	inputCollector := collector.New()

	steps := []struct {
		key         string
		message     string
		errMessage  string
		result      *string
		validator   func(string) bool
		callbackBtn *tele.Btn
	}{
		{
			key:         "advertiserId",
			message:     fmt.Sprintf("Давайте создадим кампанию! %v\n\n Введите ID рекламодателя (его можно получить в API):", emoji.EMail),
			errMessage:  "Некорректное значение ID рекламодателя, введите ещё раз:",
			result:      new(string),
			validator:   func(s string) bool { return uuid.Validate(s) == nil },
			callbackBtn: nil,
		},
		{
			key:        "impressionLimit",
			message:    "Введите лимит показов:",
			errMessage: "Некорректное значение лимита показов, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.Atoi(s)
				return err == nil && f > 0
			},
			callbackBtn: nil,
		},
		{
			key:        "clicksLimit",
			message:    "Введите лимит кликов:",
			errMessage: "Некорректное значение лимита кликов, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.Atoi(s)
				return err == nil && f > 0
			},
			callbackBtn: nil,
		},
		{
			key:        "costPerImpression",
			message:    "Введите стоимость показа:",
			errMessage: "Некорректное значение стоимости показа, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.ParseFloat(s, 64)
				return err == nil && f > 0
			},
			callbackBtn: nil,
		},
		{
			key:        "costPerClick",
			message:    "Введите стоимость клика:",
			errMessage: "Некорректное значение стоимости клика, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.ParseFloat(s, 64)
				return err == nil && f > 0
			},
			callbackBtn: nil,
		},
		{
			key:         "adTitle",
			message:     "Введите заголовок рекламы:",
			errMessage:  "Название должно быть длиной хотя бы 5 символов, введите ещё раз:",
			result:      new(string),
			validator:   func(s string) bool { return len(s) >= 5 },
			callbackBtn: nil,
		},
		{
			key:         "adText",
			message:     "Введите описание рекламы:",
			errMessage:  "Описание должно быть длиной хотя бы 10 символов, введите ещё раз:",
			result:      new(string),
			validator:   func(s string) bool { return len(s) >= 10 },
			callbackBtn: nil,
		},
		{
			key:        "startDate",
			message:    "Введите дату начала кампании (натуральное число):",
			errMessage: "Некорректное значение даты начала кампании, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.Atoi(s)
				return err == nil && f > 0
			},
			callbackBtn: nil,
		},
		{
			key:        "endDate",
			message:    "Введите дату окончания кампании (натуральное число):",
			errMessage: "Некорректное значение даты окончания кампании, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.Atoi(s)
				return err == nil && f > 0
			},
			callbackBtn: nil,
		},
		{
			key:        "gender",
			message:    "Введите пол для таргетирования: MALE, FEMALE или ALL",
			errMessage: "Некорректное значение пола, введите ещё раз:",
			result:     new(string),
			validator:  func(s string) bool { return s == "MALE" || s == "FEMALE" || s == "ALL" },
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_gender",
			},
		},
		{
			key:        "ageFrom",
			message:    "Введите нижнюю границу возраста пользователя для таргетирования (опционально):",
			errMessage: "Некорректное значение возраста, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.Atoi(s)
				return (err == nil && f > 0 && f < 120) || s == ""
			},
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_age_from",
			},
		},
		{
			key:        "ageTo",
			message:    "Введите верхнюю границу возраста пользователя для таргетирования (опционально):",
			errMessage: "Некорректное значение возраста, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.Atoi(s)
				return (err == nil && f > 0 && f < 120) || s == ""
			},
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_age_to",
			},
		},
		{
			key:        "location",
			message:    "Введите локацию пользователя для таргетирования (опционально):",
			errMessage: "Некорректное значение локации, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				return s == "" || len(s) > 1
			},
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_location",
			},
		},
	}

	for _, step := range steps {
		done := false

		keyboard := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{},
		}

		if step.callbackBtn != nil {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tele.InlineButton{*step.callbackBtn.Inline()})
		}

		_ = inputCollector.Send(tgCtx, step.message, keyboard)

		for !done {
			response, errGet := h.inputManager.Get(context.Background(), tgCtx.Sender().ID, 0, step.callbackBtn)
			if response.Message != nil {
				inputCollector.Collect(response.Message)
			}
			switch {
			case response.Canceled:
				_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
				return nil
			case errGet != nil || response.Message == nil:
				logger.Log.Errorf("failed to get input: %v", errGet)
				_ = inputCollector.Send(tgCtx, "Что-то пошло не так. Попробуем ещё раз")
			case response.Callback != nil:
				*step.result = ""
				_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})
				done = true
			case !step.validator(response.Message.Text):
				_ = inputCollector.Send(tgCtx, step.errMessage)
			case step.validator(response.Message.Text):
				*step.result = response.Message.Text
				_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})
				done = true
			}
		}
	}

	_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})

	createCampaignDTO.AdvertiserID = *steps[0].result
	createCampaignDTO.ImpressionsLimit = int32(parsing.IntMustParse(*steps[1].result))
	createCampaignDTO.ClicksLimit = int32(parsing.IntMustParse(*steps[2].result))
	createCampaignDTO.CostPerImpression = parsing.Float64MustParse(*steps[3].result)
	createCampaignDTO.CostPerClick = parsing.Float64MustParse(*steps[4].result)
	createCampaignDTO.AdTitle = *steps[5].result
	createCampaignDTO.AdText = *steps[6].result
	createCampaignDTO.StartDate = int32(parsing.IntMustParse(*steps[7].result))
	createCampaignDTO.EndDate = int32(parsing.IntMustParse(*steps[8].result))

	if *steps[9].result == "" {
		createCampaignDTO.Targeting.Gender = pointers.String("ALL")
	} else {
		createCampaignDTO.Targeting.Gender = steps[9].result
	}

	if *steps[10].result != "" {
		createCampaignDTO.Targeting.AgeFrom = pointers.Int32(int32(parsing.IntMustParse(*steps[10].result)))
	} else {
		createCampaignDTO.Targeting.AgeFrom = pointers.Int32(0)
	}
	if *steps[11].result != "" {
		createCampaignDTO.Targeting.AgeTo = pointers.Int32(int32(parsing.IntMustParse(*steps[11].result)))
	} else {
		createCampaignDTO.Targeting.AgeTo = pointers.Int32(999)
	}

	createCampaignDTO.Targeting.Location = steps[12].result

	result, err := h.service.CreateCampaign(traceCtx, createCampaignDTO)

	if err != nil {
		logger.Log.Errorf("failed to create campaign: %v", err)
		span.RecordError(err)
		_ = inputCollector.Send(tgCtx, fmt.Sprintf("Что-то пошло не так. Ошибка: %v", err))
	}

	_ = inputCollector.Send(tgCtx, fmt.Sprintf(`
		Кампания успешно создана! %v
		
		%v ID кампании: %v
		%v ID рекламодателя: %v
		%v Лимит показов: %v
		%v Лимит кликов: %v
		%v Стоимость показа: %v
		%v Стоимость клика: %v
		%v Заголовок рекламы: %v
		%v Текст рекламы: %v
		%v Дата начала: %v
		%v Дата окончания: %v
		%v Пол: %v
		%v Возраст от: %v
		%v Возраст до: %v
		%v Локация: %v
		%v Кампания прошла модерацию: %v
		`, emoji.ConfettiBall,
		emoji.BlackSmallSquare, result.CampaignID,
		emoji.BlackSmallSquare, result.AdvertiserID,
		emoji.BlackSmallSquare, result.ImpressionsLimit,
		emoji.BlackSmallSquare, result.ClicksLimit,
		emoji.BlackSmallSquare, result.CostPerImpression,
		emoji.BlackSmallSquare, result.CostPerClick,
		emoji.BlackSmallSquare, result.AdTitle,
		emoji.BlackSmallSquare, result.AdText,
		emoji.BlackSmallSquare, result.StartDate,
		emoji.BlackSmallSquare, result.EndDate,
		emoji.BlackSmallSquare, result.Targeting.Gender,
		emoji.BlackSmallSquare, result.Targeting.AgeFrom,
		emoji.BlackSmallSquare, result.Targeting.AgeTo,
		emoji.BlackSmallSquare, result.Targeting.Location,
		emoji.BlackSmallSquare, result.Approved))

	return nil
}

// Поскольку бот создан для инвесторов в демонстративных целях, пагинация пока что не реализована

func (h *CampaignHandler) GetAll(traceCtx context.Context, tgCtx tele.Context) error {
	tracer := otel.Tracer("tg-get-with-pagination-handler")
	traceCtx, span := tracer.Start(context.Background(), "Tg-GetWithPagination")
	defer span.End()

	var paginationDTO dto.GetCampaignsWithPaginationDTO

	paginationDTO.Limit = 99999
	paginationDTO.Offset = 0

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
			paginationDTO.AdvertiserID = response.Message.Text
			_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})
			done = true
		}
	}

	result, err := h.service.GetCampaignWithPagination(traceCtx, paginationDTO)

	if err != nil {
		logger.Log.Errorf("failed to create campaign: %v", err)
		span.RecordError(err)
		_ = inputCollector.Send(tgCtx, fmt.Sprintf("Что-то пошло не так. Ошибка: %v", err))
	}

	_ = inputCollector.Send(tgCtx, "Ваши кампании: ")

	for _, v := range result {
		_ = inputCollector.Send(tgCtx, fmt.Sprintf(`
									%v ID кампании: %v
									%v ID рекламодателя: %v
									%v Лимит показов: %v
									%v Лимит кликов: %v
									%v Стоимость показа: %v
									%v Стоимость клика: %v
									%v Заголовок рекламы: %v
									%v Текст рекламы: %v
									%v Дата начала: %v
									%v Дата окончания: %v
									%v Пол: %v
									%v Возраст от: %v
									%v Возраст до: %v
									%v Локация: %v
									%v Кампания прошла модерацию: %v`,
			emoji.BlackSmallSquare, v.CampaignID,
			emoji.BlackSmallSquare, v.AdvertiserID,
			emoji.BlackSmallSquare, v.ImpressionsLimit,
			emoji.BlackSmallSquare, v.ClicksLimit,
			emoji.BlackSmallSquare, v.CostPerImpression,
			emoji.BlackSmallSquare, v.CostPerClick,
			emoji.BlackSmallSquare, v.AdTitle,
			emoji.BlackSmallSquare, v.AdText,
			emoji.BlackSmallSquare, v.StartDate,
			emoji.BlackSmallSquare, v.EndDate,
			emoji.BlackSmallSquare, v.Targeting.Gender,
			emoji.BlackSmallSquare, v.Targeting.AgeFrom,
			emoji.BlackSmallSquare, v.Targeting.AgeTo,
			emoji.BlackSmallSquare, v.Targeting.Location,
			emoji.BlackSmallSquare, v.Approved))
	}

	return nil
}

func (h *CampaignHandler) UpdateCampaign(traceCtx context.Context, tgCtx tele.Context) error {
	tracer := otel.Tracer("tg-update-campaign-handler")
	traceCtx, span := tracer.Start(context.Background(), "Tg-UpdateCampaign")
	defer span.End()

	var updateDTO dto.UpdateCampaignDTO

	inputCollector := collector.New()

	steps := []struct {
		key         string
		message     string
		errMessage  string
		result      *string
		validator   func(string) bool
		callbackBtn *tele.Btn
	}{
		{
			key:         "advertiserId",
			message:     fmt.Sprintf("Давайте обновим кампанию! %v\n\n Введите ID рекламодателя (его можно получить в API):", emoji.EMail),
			errMessage:  "Некорректное значение ID рекламодателя, введите ещё раз:",
			result:      new(string),
			validator:   func(s string) bool { return uuid.Validate(s) == nil },
			callbackBtn: nil,
		},
		{
			key:         "campaignId",
			message:     "Введите ID кампании (его можно получить в API):",
			errMessage:  "Некорректное значение ID кампании, введите ещё раз:",
			result:      new(string),
			validator:   func(s string) bool { return uuid.Validate(s) == nil },
			callbackBtn: nil,
		},
		{
			key:        "impressionLimit",
			message:    "Введите новый лимит показов:",
			errMessage: "Некорректное значение лимита показов, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.Atoi(s)
				return err == nil && f > 0
			},
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_impression_limit",
			},
		},
		{
			key:        "clicksLimit",
			message:    "Введите новый лимит кликов:",
			errMessage: "Некорректное значение лимита кликов, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.Atoi(s)
				return err == nil && f > 0
			},
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_clicks_limit",
			},
		},
		{
			key:        "costPerImpression",
			message:    "Введите новую стоимость показа:",
			errMessage: "Некорректное значение стоимости показа, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.ParseFloat(s, 64)
				return err == nil && f > 0
			},
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_cost_per_impression",
			},
		},
		{
			key:        "costPerClick",
			message:    "Введите новую стоимость клика:",
			errMessage: "Некорректное значение стоимости клика, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.ParseFloat(s, 64)
				return err == nil && f > 0
			},
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_cost_per_click",
			},
		},
		{
			key:        "adTitle",
			message:    "Введите новый заголовок рекламы:",
			errMessage: "Название должно быть длиной хотя бы 5 символов, введите ещё раз:",
			result:     new(string),
			validator:  func(s string) bool { return len(s) >= 5 },
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_ad_title",
			},
		},
		{
			key:        "adText",
			message:    "Введите новое описание рекламы:",
			errMessage: "Описание должно быть длиной хотя бы 10 символов, введите ещё раз:",
			result:     new(string),
			validator:  func(s string) bool { return len(s) >= 10 },
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_ad_text",
			},
		},
		{
			key:        "gender",
			message:    "Введите новый пол для таргетирования: MALE, FEMALE или ALL",
			errMessage: "Некорректное значение пола, введите ещё раз:",
			result:     new(string),
			validator:  func(s string) bool { return s == "MALE" || s == "FEMALE" || s == "ALL" },
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_gender",
			},
		},
		{
			key:        "ageFrom",
			message:    "Введите новую нижнюю границу возраста пользователя для таргетирования:",
			errMessage: "Некорректное значение возраста, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.Atoi(s)
				return (err == nil && f > 0 && f < 120) || s == ""
			},
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_age_from",
			},
		},
		{
			key:        "ageTo",
			message:    "Введите новую верхнюю границу возраста пользователя для таргетирования:",
			errMessage: "Некорректное значение возраста, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				f, err := strconv.Atoi(s)
				return (err == nil && f > 0 && f < 120) || s == ""
			},
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_age_to",
			},
		},
		{
			key:        "location",
			message:    "Введите новую локацию пользователя для таргетирования:",
			errMessage: "Некорректное значение локации, введите ещё раз:",
			result:     new(string),
			validator: func(s string) bool {
				return s == "" || len(s) > 1
			},
			callbackBtn: &tele.Btn{
				Text:   fmt.Sprintf("Пропустить %v", emoji.RightArrow),
				Unique: "skip_location",
			},
		},
	}

	for _, step := range steps {
		done := false

		keyboard := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{},
		}

		if step.callbackBtn != nil {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tele.InlineButton{*step.callbackBtn.Inline()})
		}

		_ = inputCollector.Send(tgCtx, step.message, keyboard)

		for !done {
			response, errGet := h.inputManager.Get(context.Background(), tgCtx.Sender().ID, 0, step.callbackBtn)
			if response.Message != nil {
				inputCollector.Collect(response.Message)
			}
			switch {
			case response.Canceled:
				_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
				return nil
			case errGet != nil || response.Message == nil:
				logger.Log.Errorf("failed to get input: %v", errGet)
				_ = inputCollector.Send(tgCtx, "Что-то пошло не так. Попробуем ещё раз")
			case response.Callback != nil:
				*step.result = ""
				_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})
				done = true
			case !step.validator(response.Message.Text):
				_ = inputCollector.Send(tgCtx, step.errMessage)
			case step.validator(response.Message.Text):
				*step.result = response.Message.Text
				_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})
				done = true
			}
		}
	}

	_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})

	updateDTO.AdvertiserID = *steps[0].result
	updateDTO.CampaignID = *steps[1].result
	updateDTO.ImpressionsLimit = parsing.Int32PointerMustParse(*steps[2].result)
	updateDTO.ClicksLimit = parsing.Int32PointerMustParse(*steps[3].result)
	updateDTO.CostPerImpression = parsing.Float64MustParse(*steps[4].result)
	updateDTO.CostPerClick = parsing.Float64MustParse(*steps[5].result)
	updateDTO.AdTitle = *steps[6].result
	updateDTO.AdText = *steps[7].result

	updateDTO.Targeting.Gender = steps[8].result

	if *steps[9].result != "" {
		updateDTO.Targeting.AgeFrom = int32(parsing.IntMustParse(*steps[9].result))
	} else {
		updateDTO.Targeting.AgeFrom = 0
	}
	if *steps[10].result != "" {
		updateDTO.Targeting.AgeTo = int32(parsing.IntMustParse(*steps[10].result))
	} else {
		updateDTO.Targeting.AgeTo = 999
	}
	updateDTO.Targeting.Location = *steps[11].result

	result, err := h.service.UpdateCampaign(traceCtx, updateDTO)

	if err != nil {
		span.RecordError(err)
		_ = inputCollector.Send(tgCtx, fmt.Sprintf("Что-то пошло не так: %v", err))
		return err
	}

	_ = inputCollector.Send(tgCtx, fmt.Sprintf(`
		Кампания успешно обновлена! %v
		
		%v ID кампании: %v
		%v ID рекламодателя: %v
		%v Лимит показов: %v
		%v Лимит кликов: %v
		%v Стоимость показа: %v
		%v Стоимость клика: %v
		%v Заголовок рекламы: %v
		%v Текст рекламы: %v
		%v Дата начала: %v
		%v Дата окончания: %v
		%v Пол: %v
		%v Возраст от: %v
		%v Возраст до: %v
		%v Локация: %v
		%v Кампания прошла модерацию: %v
		`, emoji.ConfettiBall,
		emoji.BlackSmallSquare, result.CampaignID,
		emoji.BlackSmallSquare, result.AdvertiserID,
		emoji.BlackSmallSquare, result.ImpressionsLimit,
		emoji.BlackSmallSquare, result.ClicksLimit,
		emoji.BlackSmallSquare, result.CostPerImpression,
		emoji.BlackSmallSquare, result.CostPerClick,
		emoji.BlackSmallSquare, result.AdTitle,
		emoji.BlackSmallSquare, result.AdText,
		emoji.BlackSmallSquare, result.StartDate,
		emoji.BlackSmallSquare, result.EndDate,
		emoji.BlackSmallSquare, result.Targeting.Gender,
		emoji.BlackSmallSquare, result.Targeting.AgeFrom,
		emoji.BlackSmallSquare, result.Targeting.AgeTo,
		emoji.BlackSmallSquare, result.Targeting.Location,
		emoji.BlackSmallSquare, result.Approved))

	return nil
}

func (h *CampaignHandler) DeleteCampaign(traceCtx context.Context, tgCtx tele.Context) error {
	tracer := otel.Tracer("tg-delete-handler")
	traceCtx, span := tracer.Start(context.Background(), "Tg-DeleteCampaign")
	defer span.End()

	var campaignDTO dto.DeleteCampaignDTO

	inputCollector := collector.New()

	steps := []struct {
		key         string
		message     string
		errMessage  string
		result      *string
		validator   func(string) bool
		callbackBtn *tele.Btn
	}{
		{
			key:         "advertiserId",
			message:     fmt.Sprintf("Давайте обновим кампанию! %v\n\n Введите ID рекламодателя (его можно получить в API):", emoji.EMail),
			errMessage:  "Некорректное значение ID рекламодателя, введите ещё раз:",
			result:      new(string),
			validator:   func(s string) bool { return uuid.Validate(s) == nil },
			callbackBtn: nil,
		},
		{
			key:         "campaignId",
			message:     "Введите ID кампании (его можно получить в API):",
			errMessage:  "Некорректное значение ID кампании, введите ещё раз:",
			result:      new(string),
			validator:   func(s string) bool { return uuid.Validate(s) == nil },
			callbackBtn: nil,
		},
	}

	for _, step := range steps {
		done := false

		keyboard := &tele.ReplyMarkup{
			InlineKeyboard: [][]tele.InlineButton{},
		}

		if step.callbackBtn != nil {
			keyboard.InlineKeyboard = append(keyboard.InlineKeyboard, []tele.InlineButton{*step.callbackBtn.Inline()})
		}

		_ = inputCollector.Send(tgCtx, step.message, keyboard)

		for !done {
			response, errGet := h.inputManager.Get(context.Background(), tgCtx.Sender().ID, 0, step.callbackBtn)
			if response.Message != nil {
				inputCollector.Collect(response.Message)
			}
			switch {
			case response.Canceled:
				_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true, ExcludeLast: true})
				return nil
			case errGet != nil || response.Message == nil:
				logger.Log.Errorf("failed to get input: %v", errGet)
				_ = inputCollector.Send(tgCtx, "Что-то пошло не так. Попробуем ещё раз")
			case response.Callback != nil:
				*step.result = ""
				_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})
				done = true
			case !step.validator(response.Message.Text):
				_ = inputCollector.Send(tgCtx, step.errMessage)
			case step.validator(response.Message.Text):
				*step.result = response.Message.Text
				_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})
				done = true
			}
		}
	}

	_ = inputCollector.Clear(tgCtx, collector.ClearOptions{IgnoreErrors: true})

	campaignDTO.AdvertiserID = *steps[0].result
	campaignDTO.CampaignID = *steps[1].result

	err := h.service.DeleteCampaign(traceCtx, campaignDTO)
	if err != nil {
		logger.Log.Error(err)
		span.RecordError(err)
		_ = inputCollector.Send(tgCtx, fmt.Sprintf("Что-то пошло не так: %v", err))
		return err
	}

	_ = inputCollector.Send(tgCtx, "Кампания удалена")

	return nil
}
