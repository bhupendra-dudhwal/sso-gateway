package ports

import (
	"context"

	"go.uber.org/zap"
)

type Logger interface {
	With(fields ...zap.Field) Logger
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	InfoCtx(ctx context.Context, msg string, fields ...zap.Field)
	ErrorCtx(ctx context.Context, msg string, fields ...zap.Field)
	WarnCtx(ctx context.Context, msg string, fields ...zap.Field)
	DebugCtx(ctx context.Context, msg string, fields ...zap.Field)
}
