package wxutil

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/medivhzhan/weapp/v3"
	"github.com/redis/go-redis/v9"
)

func InitMiniApp(redis *redis.Client, appID string) *weapp.Client {
	appSecret, _ := redis.Get(context.Background(), fmt.Sprintf("wa:secret:%s", appID)).Result()
	return weapp.NewClient(
		appID,
		appSecret,
		weapp.WithHttpClient(&http.Client{Timeout: 10 * time.Second}),
		weapp.WithAccessTokenSetter(func(appid, secret string) (token string, expireIn uint) {
			key := redis.Get(context.Background(), fmt.Sprintf("wa:access_token:%s", appID))
			str, _ := key.Result()
			return str, 10
		}),
	)
}
