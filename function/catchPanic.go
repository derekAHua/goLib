package function

import (
	"fmt"
	"github.com/derekAHua/goLib/env"
	"github.com/derekAHua/goLib/utils"
	"github.com/derekAHua/goLib/zlog"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// @Author: Derek
// @Description: Catch Panic.
// @Date: 2022/4/30 23:55
// @Version 1.0

func CatchPanic(c *gin.Context, cleanups ...func()) {
	if r := recover(); r != nil {
		fields := []zlog.Field{
			zap.String("logId", zlog.GetLogId(c)),
			zap.String("requestId", zlog.GetRequestId(c)),
			zap.String("module", env.GetAppName()),
			zap.Stack("stack"),
		}

		if c.Request != nil {
			path := c.Request.URL.Path
			raw := c.Request.URL.RawQuery
			if raw != "" {
				path = fmt.Sprintf("%s?%s", path, raw)
			}
			fields = append(fields,
				zap.String("url", path),
				zap.String("refer", c.Request.Referer()),
				zap.String("host", c.Request.Host),
				zap.String("ua", c.Request.UserAgent()),
				zap.String("clientIp", utils.GetClientIp(c)),
			)
		}

		zlog.ErrorLogger(c, zlog.LogNameServer, fmt.Sprintf("Panic[recover](%v)", r), fields...)

		for _, fn := range cleanups {
			fn()
		}
	}
}
