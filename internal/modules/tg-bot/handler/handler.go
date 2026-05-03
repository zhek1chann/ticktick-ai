package handler

import (
	tgbot "ticktick-ai/internal/modules/tg-bot/service"

	tele "gopkg.in/telebot.v3"
)

type Handler struct {
	svc *tgbot.BotService
}

func New(svc *tgbot.BotService) *Handler {
	return &Handler{svc: svc}
}

func (h *Handler) RegisterRoutes(bot *tele.Bot) {
	bot.Handle("/start", h.StartHandler)
	bot.Handle(tele.OnText, h.OnTextHandler)
	bot.Handle(tele.OnVoice, h.OnVoiceHandler)
}
