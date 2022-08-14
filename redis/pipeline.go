package redis

import (
	"errors"
	"github.com/derekAHua/goLib/utils"
	"github.com/derekAHua/goLib/zlog"
	"go.uber.org/zap"
	"time"

	"github.com/gin-gonic/gin"
)

type PipeLiner interface {
	Exec(ctx *gin.Context) ([]interface{}, error)
	Put(ctx *gin.Context, cmd string, args ...interface{}) error
}

type commands struct {
	cmd   string
	args  []interface{}
	reply interface{}
	err   error
}

type Pipeline struct {
	cmdS  []commands
	err   error
	redis *Redis
}

func (r *Redis) Pipeline() PipeLiner {
	return &Pipeline{
		redis: r,
	}
}

func (p *Pipeline) Put(ctx *gin.Context, cmd string, args ...interface{}) error {
	if len(args) < 1 {
		return errors.New("no key found in args")
	}
	c := commands{
		cmd:  cmd,
		args: args,
	}
	p.cmdS = append(p.cmdS, c)
	return nil
}

func (p *Pipeline) Exec(ctx *gin.Context) (res []interface{}, err error) {
	start := time.Now()

	conn := p.redis.pool.Get()
	defer func() { _ = conn.Close() }()

	for i := range p.cmdS {
		err = conn.Send(p.cmdS[i].cmd, p.cmdS[i].args...)
	}

	err = conn.Flush()

	var msg string
	var ralCode int
	if err == nil {
		ralCode = 0
		for i := range p.cmdS {
			var reply interface{}
			reply, err = conn.Receive()
			res = append(res, reply)
			p.cmdS[i].reply, p.cmdS[i].err = reply, err
		}

		msg = "pipeline exec success"
	} else {
		ralCode = -1
		p.err = err
		msg = "pipeline exec error: " + err.Error()
	}

	end := time.Now()
	fields := []zlog.Field{
		zlog.WithTopicField(zlog.LogNameRedis),
		zap.String("protobuf", "redis"),
		zap.String("remoteAddr", p.redis.RemoteAddr),
		zap.String("service", p.redis.Service),
		zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
		zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zap.Float64("cost", utils.GetRequestCost(start, end)),
		zap.Int("ralCode", ralCode),
	}

	zlog.InfoLogger(ctx, zlog.LogNameRedis, msg, fields...)

	return res, err
}
