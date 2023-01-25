package log

import (
	"context"

	"go.uber.org/zap"
)

type ctxKey int

const (
	loggerKey ctxKey = iota
)

var noopLogger = zap.NewNop()

func WithLogger(ctx context.Context, set *zap.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, set)
}

func GetLogger(ctx context.Context) *zap.Logger {
	if ctx == nil {
		return noopLogger
	}
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger
	}
	return noopLogger
}
