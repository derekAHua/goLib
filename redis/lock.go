package redis

import (
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/pkg/errors"
)

const (
	EXSeconds       = "EX"
	PXMillisSeconds = "PX"
	NOTExists       = "NX"
)

// SetNxByEX 设置过期时间为秒级的redis分布式锁
func (r *Redis) SetNxByEX(ctx *gin.Context, key string, value interface{}, expire uint64) (bool, error) {
	return r.tryLock(ctx, key, value, expire, EXSeconds)
}

// SetNxByPX 设置过期时间为毫秒的redis分布式锁
func (r *Redis) SetNxByPX(ctx *gin.Context, key string, value interface{}, expire uint64) (bool, error) {
	return r.tryLock(ctx, key, value, expire, PXMillisSeconds)
}

func (r *Redis) tryLock(ctx *gin.Context, key string, value interface{}, expire uint64, exType string) (bool, error) {
	str := parseToString(value)
	if str == "" {
		return false, errors.New("value is empty")
	}

	_, err := redis.String(r.Do(ctx, "SET", key, str, exType, expire, NOTExists))

	if err == redis.ErrNil {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
