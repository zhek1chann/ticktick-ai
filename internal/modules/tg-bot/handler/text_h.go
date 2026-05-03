package handler

import (
	"context"
	"time"

	tele "gopkg.in/telebot.v3"
)

func (h *Handler) OnTextHandler(c tele.Context) error {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	return c.Send(h.svc.ProcessText(ctx, c.Text()))
}
