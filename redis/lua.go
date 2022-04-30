package redis

import (
	"github.com/derekAHua/goLib/utils"
	"github.com/derekAHua/goLib/zlog"
	"go.uber.org/zap"
	"time"

	"github.com/gin-gonic/gin"
	redigo "github.com/gomodule/redigo/redis"
)

func (r *Redis) Lua(ctx *gin.Context, script string, keyCount int, keysAndArgs ...interface{}) (interface{}, error) {
	start := time.Now()

	lua := redigo.NewScript(keyCount, script)
	conn := r.pool.Get()
	defer conn.Close()

	reply, err := lua.Do(conn, keysAndArgs...)

	ralCode := 0
	msg := "pipeline exec succ"
	if err != nil {
		ralCode = -1
		msg = "pipeline exec error: " + err.Error()
	}
	end := time.Now()

	fields := []zlog.Field{
		zlog.WithTopicField(zlog.LogNameModule),
		zap.String("prot", "redis"),
		zap.String("remoteAddr", r.RemoteAddr),
		zap.String("service", r.Service),
		zap.String("requestStartTime", utils.GetFormatRequestTime(start)),
		zap.String("requestEndTime", utils.GetFormatRequestTime(end)),
		zap.Float64("cost", utils.GetRequestCost(start, end)),
		zap.Int("ralCode", ralCode),
	}

	zlog.InfoLogger(ctx, zlog.LogNameAccess, msg, fields...)

	return reply, err
}
