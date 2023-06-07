package graphqlutil

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	var ok bool
	var midaErr midacode.Error2
	tmpErr := e

	for {
		tmpErr = errors.Unwrap(tmpErr)
		if tmpErr == nil {
			break
		}
		midaErr, ok = tmpErr.(midacode.Error2)
		if ok {
			break
		}
	}
	err := graphql.DefaultErrorPresenter(ctx, e)
	err.Extensions = map[string]interface{}{"cn": midaErr.CN()}

	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, midaErr.CN())
	span.RecordError(e)

	return err
}
