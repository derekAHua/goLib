package base

import (
	"github.com/derekAHua/goLib/redis"
	"github.com/gin-gonic/gin"
)

type RedisClient struct {
	*redis.Redis
}

func (r *RedisClient) Do(ctx *gin.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	return r.Redis.Do(ctx, commandName, args...)
}
