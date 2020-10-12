package log

import "context"

type (
	traceIDKey struct{}
	userIDKey  struct{}
)

func NewTraceIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, traceIDKey{}, traceID)
}

func FromTraceIDContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(traceIDKey{}).(string)
	return id, ok
}

func NewUserIDContext(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, userIDKey{}, traceID)
}

func FromUserIDContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey{}).(string)
	return id, ok
}
