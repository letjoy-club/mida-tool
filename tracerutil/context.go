package tracerutil

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

func WithSpanContext(ctx context.Context, header http.Header) context.Context {
	ctx = otel.GetTextMapPropagator().Extract(
		ctx, propagation.HeaderCarrier(header),
	)
	return ctx
}

func GetSpan(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}
