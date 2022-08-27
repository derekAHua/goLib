package zlog

import (
	"github.com/derekAHua/goLib/env"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var mapZapLogger = make(map[LogName]*zap.Logger)

func zapLogger(ctx *gin.Context, logName LogName) *zap.Logger {
	ret, ok := mapZapLogger[logName]
	if !ok {
		ret = mapZapLogger[LogNameServer]
	}

	if ctx == nil {
		return ret
	}

	return ret.With(
		zap.String("logId", GetLogId(ctx)),
		zap.String("requestId", GetRequestId(ctx)),
		zap.String("module", env.GetAppName()),
		zap.String("localIp", env.GetLocalIP()),
		zap.String("uri", ctx.GetString(ContextKeyUri)),
	)
}

func DebugLogger(ctx *gin.Context, logName LogName, msg string, fields ...zap.Field) {
	if NoLog(ctx) {
		return
	}
	zapLogger(ctx, logName).Debug(msg, fields...)
}

func InfoLogger(ctx *gin.Context, logName LogName, msg string, fields ...zap.Field) {
	if NoLog(ctx) {
		return
	}
	zapLogger(ctx, logName).Info(msg, fields...)
}

func WarnLogger(ctx *gin.Context, logName LogName, msg string, fields ...zap.Field) {
	if NoLog(ctx) {
		return
	}
	zapLogger(ctx, logName).Warn(msg, fields...)
}

func ErrorLogger(ctx *gin.Context, logName LogName, msg string, fields ...zap.Field) {
	if NoLog(ctx) {
		return
	}
	zapLogger(ctx, logName).Error(msg, fields...)
}

func PanicLogger(ctx *gin.Context, logName LogName, msg string, fields ...zap.Field) {
	if NoLog(ctx) {
		return
	}
	zapLogger(ctx, logName).Panic(msg, fields...)
}

func FatalLogger(ctx *gin.Context, logName LogName, msg string, fields ...zap.Field) {
	if NoLog(ctx) {
		return
	}
	zapLogger(ctx, logName).Fatal(msg, fields...)
}
