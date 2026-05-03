package logger

import (
	"context"
	"log/slog"
	"ticktick-ai/pkg/db"
	"ticktick-ai/pkg/db/prettier"
)

func LogQuery(ctx context.Context, q db.Query, args ...interface{}) {
	prettyQuery := prettier.Pretty(q.QueryRaw, prettier.PlaceholderDollar, args...)
	slog.InfoContext(ctx, "", "sql", q.Name,
		slog.String("query", prettyQuery),
	)
}
