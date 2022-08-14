package redis

import (
	"fmt"
	"github.com/derekAHua/goLib/utils"
	"github.com/derekAHua/goLib/zlog"
	"go.uber.org/zap"
	"time"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
)

// 日志打印Do args部分支持的最大长度
const logForRedisValue = 50

type Conf struct {
	Service         string        `yaml:"service"`
	Addr            string        `yaml:"addr"`
	Password        string        `yaml:"password"`
	MaxIdle         int           `yaml:"maxIdle"`
	MaxActive       int           `yaml:"maxActive"`
	IdleTimeout     time.Duration `yaml:"idleTimeout"`
	MaxConnLifetime time.Duration `yaml:"maxConnLifetime"`
	ConnTimeOut     time.Duration `yaml:"connTimeOut"`
	ReadTimeOut     time.Duration `yaml:"readTimeOut"`
	WriteTimeOut    time.Duration `yaml:"writeTimeOut"`
}

func (conf *Conf) checkConf() {
	if conf.MaxIdle == 0 {
		conf.MaxIdle = 50
	}
	if conf.MaxActive == 0 {
		conf.MaxActive = 100
	}
	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = 3 * time.Minute
	}
	if conf.MaxConnLifetime == 0 {
		conf.MaxConnLifetime = 10 * time.Minute
	}
	if conf.ConnTimeOut == 0 {
		conf.ConnTimeOut = 1200 * time.Millisecond
	}
	if conf.ReadTimeOut == 0 {
		conf.ReadTimeOut = 1200 * time.Millisecond
	}
	if conf.WriteTimeOut == 0 {
		conf.WriteTimeOut = 1200 * time.Millisecond
	}
}

type Redis struct {
	pool       *redigo.Pool
	Service    string
	RemoteAddr string
}

func InitRedisClient(conf Conf) (*Redis, error) {
	conf.checkConf()
	p := &redigo.Pool{
		MaxIdle:         conf.MaxIdle,
		MaxActive:       conf.MaxActive,
		IdleTimeout:     conf.IdleTimeout,
		MaxConnLifetime: conf.MaxConnLifetime,
		Wait:            true,
		Dial: func() (conn redigo.Conn, e error) {
			con, err := redigo.Dial(
				"tcp",
				conf.Addr,
				redigo.DialPassword(conf.Password),
				redigo.DialConnectTimeout(conf.ConnTimeOut),
				redigo.DialReadTimeout(conf.ReadTimeOut),
				redigo.DialWriteTimeout(conf.WriteTimeOut),
			)
			if err != nil {
				return nil, err
			}
			return con, nil
		},
		TestOnBorrow: func(c redigo.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
	c := &Redis{
		Service:    conf.Service,
		RemoteAddr: conf.Addr,
		pool:       p,
	}
	return c, nil
}

func (r *Redis) Do(ctx *gin.Context, commandName string, args ...interface{}) (reply interface{}, err error) {
	start := time.Now()

	conn := r.pool.Get()
	err = conn.Err()
	if err != nil {
		zlog.ErrorLogger(ctx, zlog.LogNameRedis, "get connection error: "+err.Error(), zlog.WithTopicField(zlog.LogNameRedis), zap.String("protobuf", "redis"))
		return reply, err
	}

	reply, err = conn.Do(commandName, args...)
	if e := conn.Close(); e != nil {
		zlog.WarnLogger(ctx, zlog.LogNameRedis, "connection close error: "+e.Error(), zlog.WithTopicField(zlog.LogNameRedis), zap.String("protobuf", "redis"))
	}

	end := time.Now()

	// 执行时间 单位:毫秒
	ralCode := 0
	msg := "redis do success"
	if err != nil {
		ralCode = -1
		msg = fmt.Sprintf("redis do error: %s", err.Error())
		zlog.ErrorLogger(ctx, zlog.LogNameRedis, msg, zlog.WithTopicField(zlog.LogNameRedis), zap.String("protobuf", "redis"))
	}

	fields := []zlog.Field{
		zlog.WithTopicField(zlog.LogNameRedis),
		zap.String("protobuf", "redis"),
		zap.String("service", r.Service),
		zap.String("remoteAddr", r.RemoteAddr),
		zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
		zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zap.Float64("cost", utils.GetRequestCost(start, end)),
		zap.String("command", commandName),
		zap.String("commandVal", utils.JoinArgs(logForRedisValue, args)),
		zap.Int("ralCode", ralCode),
	}

	zlog.InfoLogger(ctx, zlog.LogNameRedis, msg, fields...)
	return reply, err
}

func (r *Redis) Close() error {
	return r.pool.Close()
}

func (r *Redis) Stats() (inUseCount, idleCount, activeCount int) {
	stats := r.pool.Stats()
	idleCount = stats.IdleCount
	activeCount = stats.ActiveCount
	inUseCount = activeCount - idleCount
	return inUseCount, idleCount, activeCount
}
