package otel

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

func Trace(ctx context.Context, spanName string) (context.Context, trace.Span) {
	tracer := otel.Tracer("medical-service")
	return tracer.Start(ctx, spanName)
}
