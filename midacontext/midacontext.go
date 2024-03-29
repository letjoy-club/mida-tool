package midacontext

import (
	"context"
	"net/http"
	"time"

	"github.com/Khan/genqlient/graphql"
	"github.com/go-chi/cors"
	"github.com/letjoy-club/mida-tool/authenticator"
	"github.com/letjoy-club/mida-tool/clienttoken"
	"github.com/medivhzhan/weapp/v3"
)

type startTime struct{}

func WithStartTime(ctx context.Context) context.Context {
	return context.WithValue(ctx, startTime{}, time.Now())
}

func GetStartTime(ctx context.Context) time.Time {
	return ctx.Value(startTime{}).(time.Time)
}

/**
 * Wechat
 */
type wechatKey struct{}

type WechatConf struct {
	Env         string
	AppID       string
	AppSecret   string
	WeappClient *weapp.Client
}

func WithWechatConf(ctx context.Context, conf WechatConf) context.Context {
	return context.WithValue(ctx, wechatKey{}, conf)
}

func GetWechatConf(ctx context.Context) WechatConf {
	return ctx.Value(wechatKey{}).(WechatConf)
}

/**
 * Tencent Map
 */
type mapKey struct{}

func WithMapConf(ctx context.Context, conf string) context.Context {
	return context.WithValue(ctx, mapKey{}, conf)
}

func GetMapConf(ctx context.Context) string {
	return ctx.Value(mapKey{}).(string)
}

/**
 * Loader
 */
type loaderKey struct{}

// loader 根据不同的业务场景，可以是不同的类型
func WithLoader[LoaderType any](ctx context.Context, loader LoaderType) context.Context {
	return context.WithValue(ctx, loaderKey{}, loader)
}

// loader 根据不同的业务场景，可以是不同的类型
func GetLoader[LoaderType any](ctx context.Context) *LoaderType {
	return ctx.Value(loaderKey{}).(*LoaderType)
}

var WithCORS = cors.Handler(cors.Options{
	// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
	AllowedOrigins:   []string{"https://*", "http://*"},
	AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
	AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "x-mita-token", "x-mida-token"},
	ExposedHeaders:   []string{"Link", "Content-Type"},
	AllowCredentials: false,
	MaxAge:           300, // Maximum value not ignored by any of major browsers
})

/**
 * Authenticator
 */
type authKey struct{}

func WithAuthenticator(ctx context.Context, auth authenticator.Authenticator) context.Context {
	return context.WithValue(ctx, authKey{}, auth)
}

func GetAuthenticator(ctx context.Context) authenticator.Authenticator {
	return ctx.Value(authKey{}).(authenticator.Authenticator)
}

/**
 * ClientToken
 */
type clientTokenKey struct{}

func WithClientToken(ctx context.Context, token clienttoken.ClientToken) context.Context {
	return context.WithValue(ctx, clientTokenKey{}, token)
}

func GetClientToken(ctx context.Context) clienttoken.ClientToken {
	return ctx.Value(clientTokenKey{}).(clienttoken.ClientToken)
}

type servicesKey struct{}

type Services struct {
	// 基础服务
	Hoopoe graphql.Client
	// IM 服务
	Smew graphql.Client
	// 匹配服务
	Whale graphql.Client
	// 通知服务
	Scream graphql.Client
}

func WithServices(ctx context.Context, services Services) context.Context {
	return context.WithValue(ctx, servicesKey{}, services)
}

func GetServices(ctx context.Context) Services {
	return ctx.Value(servicesKey{}).(Services)
}

func NewServices(url, token string) graphql.Client {
	client := newHttpClient(token)
	return graphql.NewClient(url, client)
}

/**
 * GraphQL
 */
type GraphQLResp struct {
	Data   interface{}  `json:"data"`
	Errors []GraphQLErr `json:"errors"`
}

type GraphQLErr struct {
	Message string   `json:"message"`
	Path    []string `json:"path"`
}

func ParseToken(r *http.Request, auth authenticator.Authenticator) clienttoken.ClientToken {
	var tokenStr string
	id := r.Header.Get("X-Mida-Id")
	if id != "" {
		return clienttoken.ClientToken(id)
	}
	token := r.Header.Get("X-Mida-Token")
	if token != "" {
		var err error
		tokenStr, err = auth.Verify(token)
		if err != nil {
			return clienttoken.ClientToken("$invalid: " + err.Error())
		}
	}
	return clienttoken.ClientToken(tokenStr)
}
