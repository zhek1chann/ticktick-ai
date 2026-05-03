package handler

import (
	"context"
	"io"
	"log/slog"
	"time"

	tele "gopkg.in/telebot.v3"
)

func (h *Handler) OnVoiceHandler(c tele.Context) error {
	msg := c.Message()
	if msg == nil || msg.Voice == nil {
		return c.Send("Не нашел голосовое сообщение. Попробуй еще раз.")
	}

	reader, err := c.Bot().File(&msg.Voice.File)
	if err != nil {
		slog.Error("download voice failed", "err", err)
		return c.Send("Не смог скачать голосовое сообщение. Попробуй еще раз.")
	}
	defer reader.Close()

	audio, err := io.ReadAll(reader)
	if err != nil {
		slog.Error("read voice failed", "err", err)
		return c.Send("Не смог прочитать голосовое сообщение. Попробуй еще раз.")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	return c.Send(h.svc.ProcessAudio(ctx, audio, msg.Voice.MIME))
}
