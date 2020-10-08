package log

import (
	"context"

	"go.uber.org/zap"
)

//sugar log

const (
	TraceIDKey = "trace_id"
)

func startSpan(ctx context.Context) *zap.SugaredLogger {
	if ctx == nil {
		ctx = context.Background()
	}
	sugar := logger.opts.sugaredLogger
	if id, ok := FromTraceIDContext(ctx); ok {
		sugar = sugar.With(TraceIDKey, id)
	}
	return sugar
}

func Info(ctx context.Context, args ...interface{}) {
	if logger == nil {
		return
	}
	startSpan(ctx).Info(args...)
}

func Error(ctx context.Context, args ...interface{}) {
	if logger == nil {
		return
	}
	startSpan(ctx).Error(args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	if logger == nil {
		return
	}
	startSpan(ctx).Fatal(args...)
}

func Debug(ctx context.Context, args ...interface{}) {
	if logger == nil {
		return
	}
	startSpan(ctx).Debug(args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	if logger == nil {
		return
	}
	startSpan(ctx).Panic(args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	if logger == nil {
		return
	}
	startSpan(ctx).Warn(args...)
}

func Debugf(template string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.opts.sugaredLogger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.opts.sugaredLogger.Infof(template, args...)
}

func Errorf(template string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.opts.sugaredLogger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.opts.sugaredLogger.Fatalf(template, args...)
}
