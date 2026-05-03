package handler

import tele "gopkg.in/telebot.v3"

func (h *Handler) StartHandler(c tele.Context) error {
	return c.Send("Отправь текст или голосовое сообщение, а я создам, обновлю или завершу задачу в TickTick.")
}
