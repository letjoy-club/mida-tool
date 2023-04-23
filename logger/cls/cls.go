package serving

import (
	"fmt"
	"github.com/letjoy-club/mida-tool/authenticator"
	"github.com/letjoy-club/mida-tool/midacontext"
	cls "github.com/tencentcloud/tencentcloud-cls-sdk-go"
	"io"
	"net/http"
	"runtime/debug"
	"time"
)

var fatalStr = "fatal"
var errorStr = "error"
var warnStr = "warn"
var infoStr = "info"
var debugStr = "debug"

var urlStr = "url"
var methodStr = "method"
var stackStr = "stack"
var bodyStr = "body"
var userStr = "user"
var levelStr = "level"
var paramStr = "param"
var operationStr = "operation"
var messageStr = "message"

var SendLog = func(level string) func(client *cls.AsyncProducerClient, topic, url, method, body, param, op, user, msg, stack string) {
	return func(client *cls.AsyncProducerClient, topic, url, method, param, body, op, user, msg, stack string) {
		fmt.Println("sending a log")
		now := time.Now().Unix()
		client.SendLog(topic, &cls.Log{Contents: []*cls.Log_Content{
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

var SendFatal = SendLog(fatalStr)
var SendError = SendLog(errorStr)
var SendWarn = SendLog(warnStr)
var SendInfo = SendLog(infoStr)
var SendDebug = SendLog(debugStr)

func LoggerCtx(logger *cls.AsyncProducerClient, topicID, key string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if logger != nil && topicID != "" {
			auth := authenticator.Authenticator{Key: []byte(key)}

			fn := func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if rvr := recover(); rvr != nil {
						if rvr != http.ErrAbortHandler {
							debugStack := string(debug.Stack())
							user := midacontext.ParseToken(r, auth)
							reader := r.Body
							if reader == nil {
								SendFatal(logger, topicID, r.URL.String(), r.Method, "", "", "", user.String(), "", debugStack)
							} else {
								bodyRaw, _ := io.ReadAll(reader)
								body := string(bodyRaw)
								SendFatal(logger, topicID, r.URL.String(), r.Method, body, "", "", user.String(), "", debugStack)
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
