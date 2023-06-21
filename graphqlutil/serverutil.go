package graphqlutil

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql"
	"github.com/letjoy-club/mida-tool/logger"
	"github.com/letjoy-club/mida-tool/midacode"
	"github.com/letjoy-club/mida-tool/midacontext"
	"github.com/letjoy-club/mida-tool/qcloudutil/clsutil"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func AdminOnly(ctx context.Context, obj interface{}, next graphql.Resolver) (interface{}, error) {
	token := midacontext.GetClientToken(ctx)
	if token.IsAdmin() || token.IsInternal() {
		return next(ctx)
	}
	return nil, midacode.ErrNotPermitted
}

func AroundOperations(tr trace.Tracer, service string) func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
	return func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
		oc := graphql.GetOperationContext(ctx)
		opName := "unknown"
		if oc.OperationName != "" {
			opName = oc.OperationName
		}
		ctx, span := tr.Start(ctx, service+"."+opName, trace.WithSpanKind(trace.SpanKindServer))
		defer span.End()

		token := midacontext.GetClientToken(ctx)
		opCtx := graphql.GetOperationContext(ctx)
		logger.L.Info("request body", zap.String("token", token.String()),
			zap.String("query", opCtx.RawQuery), zap.Any("param", opCtx.Variables))

		return next(ctx)
	}
}

func ErrorPresenter(ctx context.Context, e error) *gqlerror.Error {
	var ok bool
	originErr := e
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

	if !ok {
		// 有可能 db 错误
		if errors.Is(e, gorm.ErrRecordNotFound) {
			e = midacode.ErrItemNotFound
		} else {
			e = midacode.ErrUnknownError
		}
		midaErr = e.(midacode.Error2)
	}
	err := graphql.DefaultErrorPresenter(ctx, e)
	err.Extensions = midaErr.ToExtensions()

	span := trace.SpanFromContext(ctx)
	span.SetStatus(codes.Error, midaErr.CN())
	span.RecordError(originErr)

	log := clsutil.GetGraphLogger(ctx)
	switch midaErr.LogLevel() {
	case midacode.LogLevelWarn:
		log.Warn(ctx, midaErr.CN())
	case midacode.LogLevelError:
		log.Error(ctx, originErr, midaErr.CN())
	default:
		logger.L.Error("graphql request error", zap.Error(originErr))
	}

	return err
}
