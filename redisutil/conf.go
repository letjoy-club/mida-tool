package redisutil

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/redis/go-redis/v9"
)

type RedisConf struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
}

func (r RedisConf) ConnectRedis() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     r.Address,
		Password: r.Password,
		DB:       r.DB,
	})
	return client
}

type redisKey struct{}

func WithRedis(ctx context.Context, redis *redis.Client) context.Context {
	return context.WithValue(ctx, redisKey{}, redis)
}

func GetRedis(ctx context.Context) *redis.Client {
	return ctx.Value(redisKey{}).(*redis.Client)
}

func GetLocker(ctx context.Context) *redislock.Client {
	client := ctx.Value(redisKey{}).(*redis.Client)
	return redislock.New(client)
}
