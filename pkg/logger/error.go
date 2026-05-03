package logger

import (
	"context"
	"errors"
)

type errorWithLogCtx struct {
	cause error
	ctx   logCtx
}

func (e *errorWithLogCtx) Error() string { return e.cause.Error() }

func (e *errorWithLogCtx) Unwrap() error { return e.cause }

func (e *errorWithLogCtx) LogCtx() logCtx { return e.ctx }

func WrapError(ctx context.Context, err error) error {
	if err == nil {
		return nil
	}
	// Если ошибка уже имеет контекст — не делаем двойную обёртку.
	var ew *errorWithLogCtx
	if errors.As(err, &ew) {
		return err
	}

	c := logCtx{}
	if x, ok := ctx.Value(key).(logCtx); ok {
		c = x
	}
	return &errorWithLogCtx{
		cause: err,
		ctx:   c,
	}
}

// ErrorCtx переносит logCtx из ошибки в новый контекст.
// Удобно для логгеров/мидлварей, где на вход приходит только err.
func ErrorCtx(ctx context.Context, err error) context.Context {
	var ew *errorWithLogCtx
	if errors.As(err, &ew) {
		return context.WithValue(ctx, key, ew.ctx)
	}
	return ctx
}

// ExtractLogCtx даёт прямой доступ к logCtx из ошибки (не всегда нужен, но полезно иметь).
func ExtractLogCtx(err error) (logCtx, bool) {
	var ew *errorWithLogCtx
	if errors.As(err, &ew) {
		return ew.ctx, true
	}
	return logCtx{}, false
}
