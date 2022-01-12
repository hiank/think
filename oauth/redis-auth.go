package oauth

import (
	"context"

	"github.com/go-redis/redis/v8"
)

//redisAuther Auther use redis storage
type redisAuther struct {
	ctx context.Context
	cli *redis.Client
}

//NewRedisAuther new Auther with redis storage
func NewRedisAuther(ctx context.Context, cli *redis.Client) IAuther {
	return &redisAuther{ctx, cli}
}

//Auth 检查指定token是否存在与redis
func (ra *redisAuther) Auth(token string) (uid uint64, err error) {
	cmd := ra.cli.HGet(ra.ctx, "token_uid", token)
	if err = cmd.Err(); err == nil {
		uid, err = cmd.Uint64()
	}
	return
}
