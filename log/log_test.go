package log

import (
	"context"
	"testing"
)

func TestPrintLog(t *testing.T) {
	ctx := context.TODO()
	ctx = NewTraceIDContext(ctx, "1213dsd")
	Init() /**/
	Info(ctx, "xxx1")
	Info(ctx, "xxx2")
}
