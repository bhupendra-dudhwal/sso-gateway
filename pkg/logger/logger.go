package logger

import (
	"context"
	"fmt"
	"os"

	"github.com/bhupendra-dudhwal/sso-gateway/internal/constants"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/models"
	"github.com/bhupendra-dudhwal/sso-gateway/internal/core/ports"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	logger *zap.Logger
}

func NewLogger(cfg *models.Logger, env constants.Environment) ports.Logger {
	var zapCfg zap.Config

	if env == constants.Development {
		zapCfg = zap.NewDevelopmentConfig()
	} else {
		zapCfg = zap.NewProductionConfig()
	}

	// Set log level from config
	level := zapcore.InfoLevel
	if err := level.UnmarshalText([]byte(cfg.Level)); err != nil {
		fmt.Printf("invalid log level, defaulting to info: %v\n", err)
		level = zapcore.InfoLevel
	}
	zapCfg.Level = zap.NewAtomicLevelAt(level)

	logger, err := zapCfg.Build(
		zap.AddCaller(),
		zap.AddStacktrace(zapcore.ErrorLevel),
		zap.AddCallerSkip(1),
	)
	if err != nil {
		fmt.Printf("failed to initialize logger: %+v\n", err)
		os.Exit(1)
	}

	return &zapLogger{logger: logger}
}

func (l *zapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *zapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

func (l *zapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

func (l *zapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// With Context

func (l *zapLogger) InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(zapFieldsFromContext(ctx)...).Info(msg, fields...)
}

func (l *zapLogger) ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(zapFieldsFromContext(ctx)...).Error(msg, fields...)
}

func (l *zapLogger) WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(zapFieldsFromContext(ctx)...).Warn(msg, fields...)
}

func (l *zapLogger) DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	l.logger.With(zapFieldsFromContext(ctx)...).Debug(msg, fields...)
}

// --- Context Field Extractor ---
func zapFieldsFromContext(ctx context.Context) []zap.Field {
	fields := make([]zap.Field, 0)
	if reqID, ok := ctx.Value(constants.CtxRequestID).(string); ok {
		fields = append(fields, zap.String(constants.CtxRequestID.String(), reqID))
	}
	// Add more fields from context as needed
	return fields
}
