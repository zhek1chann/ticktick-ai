package logger

import (
	"context"
	"log/slog"
	"os"
)

func InitLogging(level slog.Level) {

	handler := slog.Handler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
	}))

	handler = NewHandlerMiddleware(handler)
	slog.SetDefault(slog.New(handler))
}

type HandlerMiddlware struct {
	next slog.Handler
}

func NewHandlerMiddleware(next slog.Handler) *HandlerMiddlware {

	return &HandlerMiddlware{next: next}
}

func (h *HandlerMiddlware) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *HandlerMiddlware) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {

		if c.Message != "" {
			rec.Add("message", c.Message)
		}
		if c.ShopID != 0 {
			rec.Add("shop_id", c.ShopID)
		}
	}
	return h.next.Handle(ctx, rec)
}

func WithMessage(ctx context.Context, message string) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.Message = message
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{Message: message})
}

func WithShopID(ctx context.Context, shopID int64) context.Context {
	if c, ok := ctx.Value(key).(logCtx); ok {
		c.ShopID = shopID
		return context.WithValue(ctx, key, c)
	}
	return context.WithValue(ctx, key, logCtx{ShopID: shopID})
}

func (h *HandlerMiddlware) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithAttrs(attrs)} // не забыть обернуть, но осторожно
}

func (h *HandlerMiddlware) WithGroup(name string) slog.Handler {
	return &HandlerMiddlware{next: h.next.WithGroup(name)} // не забыть обернуть, но осторожно
}
