package oauth

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type Auther interface {
	Auth(token string) (uid uint64, err error)
}

type RedisAuther struct {
	ctx      context.Context
	redisCLI *redis.Client
}

//Auth 检查指定token是否存在与redis
func (ra *RedisAuther) Auth(token string) (uint64, error) {
	cmd := ra.redisCLI.HGet(ra.ctx, "token_uid", token)
	if cmd.Err() != nil {
		return 0, cmd.Err()
	}
	return cmd.Uint64()
}
