package clsutil

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/99designs/gqlgen/graphql"
	"github.com/letjoy-club/mida-tool/authenticator"
	"github.com/letjoy-club/mida-tool/midacontext"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

var fatalStr = "fatal"
var errorStr = "error"
var warnStr = "warn"
var infoStr = "info"
var debugStr = "debug"
var spanStr = "spanId"
var traceStr = "traceId"

var serviceStr = "service"
var urlStr = "url"
var methodStr = "method"
var bodyStr = "body"
var stackStr = "stack"
var userStr = "user"
var paramStr = "param"
var levelStr = "level"
var operationStr = "operation"
var messageStr = "message"

var SendFatal = SendLog(fatalStr)
var SendError = SendLog(errorStr)
var SendWarn = SendLog(warnStr)
var SendInfo = SendLog(infoStr)
var SendDebug = SendLog(debugStr)

var SendLog = func(level string) func(client *cls.AsyncProducerClient, topic, service, url, method, body, param, op, user, msg, stack string) {
	return func(client *cls.AsyncProducerClient, topic, service, url, method, param, body, op, user, msg, stack string) {
		fmt.Println("sending a log")
		now := time.Now().Unix()
		client.SendLog(topic, &cls.Log{Contents: []*cls.Log_Content{
			{Key: &serviceStr, Value: &service},
			{Key: &urlStr, Value: &url},
			{Key: &methodStr, Value: &method},
			{Key: &bodyStr, Value: &body},
			{Key: &stackStr, Value: &stack},
			{Key: &userStr, Value: &user},
			{Key: &paramStr, Value: &param},
			{Key: &levelStr, Value: &level},
			{Key: &operationStr, Value: &op},
			{Key: &messageStr, Value: &msg},
		}, Time: &now}, nil)
	}
}

func LoggerCtx(client *cls.AsyncProducerClient, topicID, key, service string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if client != nil && topicID != "" {
			auth := authenticator.Authenticator{Key: []byte(key)}

			fn := func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if rvr := recover(); rvr != nil {
						if rvr != http.ErrAbortHandler {
							debugStack := string(debug.Stack())
							user := midacontext.ParseToken(r, auth)
							reader := r.Body
							if reader == nil {
								SendFatal(client, topicID, service, r.URL.String(), r.Method, "", "", "", user.String(), "", debugStack)
							} else {
								bodyRaw, _ := io.ReadAll(reader)
								body := string(bodyRaw)
								SendFatal(client, topicID, service, r.URL.String(), r.Method, body, "", "", user.String(), "", debugStack)
								reader.Close()
							}
						}
						panic(rvr)
					}
				}()
				next.ServeHTTP(w, r)
			}
			return http.HandlerFunc(fn)
		}
		return next
	}
}

type graphLoggerKey struct {
}

func WithGraphLogger(ctx context.Context, logger *cls.AsyncProducerClient, topicID, service string) context.Context {
	ctx = context.WithValue(ctx, graphLoggerKey{}, &Logger{
		service: service,
		cls:     logger,
		topicID: topicID,
	})
	return ctx
}

func GetGraphLogger(ctx context.Context) ILogger {
	return ctx.Value(graphLoggerKey{}).(*Logger)
}

type ILogger interface {
	Error(ctx context.Context, err error, msg string, args ...interface{})
	Info(ctx context.Context, msg string, args ...interface{})
	Warn(ctx context.Context, msg string, args ...interface{})
}

type Logger struct {
	service string
	topicID string
	cls     *cls.AsyncProducerClient
}

// Error 上报错误日志
func (l *Logger) Error(ctx context.Context, err error, msg string, args ...interface{}) {
	msgStr := fmt.Sprintf(msg, args...)

	reqCtx := graphql.GetOperationContext(ctx)
	user := midacontext.GetClientToken(ctx)
	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()

	// 记录trace
	span.SetStatus(codes.Error, msg)
	span.RecordError(err, trace.WithStackTrace(true))

	traceID := spanContext.TraceID().String()
	spanID := spanContext.SpanID().String()
	userID := user.String()

	varData, _ := json.Marshal(reqCtx.Variables)
	varStr := ""
	if varData != nil {
		varStr = string(varData)
	}

	stack := string(debug.Stack())

	now := time.Now().Unix()
	l.cls.SendLog(l.topicID, &cls.Log{Contents: []*cls.Log_Content{
		{Key: &serviceStr, Value: &l.service},
		{Key: &levelStr, Value: &errorStr},
		{Key: &userStr, Value: &userID},

		// trace 相关
		{Key: &spanStr, Value: &spanID},
		{Key: &traceStr, Value: &traceID},

		// graphql 相关
		{Key: &paramStr, Value: &varStr},
		{Key: &bodyStr, Value: &reqCtx.RawQuery},
		{Key: &operationStr, Value: &reqCtx.OperationName},

		{Key: &messageStr, Value: &msgStr},

		{Key: &stackStr, Value: &stack},
	}, Time: &now}, nil)

	fmt.Println(errorStr, now, msgStr)
}

// Info 记录 info 信息，不会记录 graphql 上下文，不会记录 stack
func (l *Logger) Info(ctx context.Context, msg string, args ...interface{}) {
	msgStr := fmt.Sprintf(msg, args...)

	reqCtx := graphql.GetOperationContext(ctx)
	user := midacontext.GetClientToken(ctx)
	userID := user.String()

	span := trace.SpanFromContext(ctx).SpanContext()
	traceID := span.TraceID().String()
	spanID := span.SpanID().String()

	now := time.Now().Unix()

	l.cls.SendLog(l.topicID, &cls.Log{Contents: []*cls.Log_Content{
		{Key: &serviceStr, Value: &l.service},
		{Key: &levelStr, Value: &infoStr},
		{Key: &userStr, Value: &userID},

		// trace 相关
		{Key: &spanStr, Value: &spanID},
		{Key: &traceStr, Value: &traceID},

		// graphql 相关
		{Key: &operationStr, Value: &reqCtx.OperationName},

		{Key: &messageStr, Value: &msgStr},
	}, Time: &now}, nil)

	fmt.Println(infoStr, now, msgStr)
}

// Warn 上报警告日志
func (l *Logger) Warn(ctx context.Context, msg string, args ...interface{}) {
	msgStr := fmt.Sprintf(msg, args...)

	reqCtx := graphql.GetOperationContext(ctx)
	user := midacontext.GetClientToken(ctx)

	span := trace.SpanFromContext(ctx)
	spanContext := span.SpanContext()
	traceID := spanContext.TraceID().String()
	spanID := spanContext.SpanID().String()
	userID := user.String()

	varData, _ := json.Marshal(reqCtx.Variables)
	varStr := ""
	if varData != nil {
		varStr = string(varData)
	}

	stack := string(debug.Stack())

	now := time.Now().Unix()
	l.cls.SendLog(l.topicID, &cls.Log{Contents: []*cls.Log_Content{
		{Key: &serviceStr, Value: &l.service},
		{Key: &levelStr, Value: &warnStr},
		{Key: &userStr, Value: &userID},

		// trace 相关
		{Key: &spanStr, Value: &spanID},
		{Key: &traceStr, Value: &traceID},

		// graphql 相关
		{Key: &paramStr, Value: &varStr},
		{Key: &bodyStr, Value: &reqCtx.RawQuery},
		{Key: &operationStr, Value: &reqCtx.OperationName},

		{Key: &messageStr, Value: &msgStr},

		{Key: &stackStr, Value: &stack},
	}, Time: &now}, nil)

	fmt.Println(warnStr, now, msgStr)
}
