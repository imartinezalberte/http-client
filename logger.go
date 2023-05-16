package httpclient

import (
	"context"

	"go.uber.org/zap"
)

type ILogger interface {
	Debug(context.Context, string, ...zap.Field)
	Info(context.Context, string, ...zap.Field)
	Warn(context.Context, string, ...zap.Field)
	Error(context.Context, string, ...zap.Field)
	DPanic(context.Context, string, ...zap.Field)
	Panic(context.Context, string, ...zap.Field)
	Fatal(context.Context, string, ...zap.Field)
}

type logger struct {
	l *zap.Logger
}

func NewLogger(l *zap.Logger, defaultFields ...zap.Field) ILogger {
	return &logger{
		l: l.With(defaultFields...),
	}
}

// TODO: we should read the properties from a configuration file.
// Imagine a case where an integration needs more than one property to be logged.
func readCtx(ctx context.Context) zap.Field {
	return zap.Any("key", "nothing")
}

func (logger *logger) Debug(ctx context.Context, msg string, fields ...zap.Field) {
	logger.l.Debug(msg, append(fields, readCtx(ctx))...)
}

func (logger *logger) Info(ctx context.Context, msg string, fields ...zap.Field) {
	logger.l.Info(msg, append(fields, readCtx(ctx))...)
}

func (logger *logger) Warn(ctx context.Context, msg string, fields ...zap.Field) {
	logger.l.Warn(msg, append(fields, readCtx(ctx))...)
}

func (logger *logger) Error(ctx context.Context, msg string, fields ...zap.Field) {
	logger.l.Error(msg, append(fields, readCtx(ctx))...)
}

func (logger *logger) DPanic(ctx context.Context, msg string, fields ...zap.Field) {
	logger.l.DPanic(msg, append(fields, readCtx(ctx))...)
}

func (logger *logger) Panic(ctx context.Context, msg string, fields ...zap.Field) {
	logger.l.Panic(msg, append(fields, readCtx(ctx))...)
}

func (logger *logger) Fatal(ctx context.Context, msg string, fields ...zap.Field) {
	logger.l.Fatal(msg, append(fields, readCtx(ctx))...)
}
