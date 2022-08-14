package zlog

import (
	"github.com/derekAHua/goLib/env"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var sugaredLogger *zap.SugaredLogger

func sLogger(ctx *gin.Context) *zap.SugaredLogger {
	if ctx == nil {
		return sugaredLogger
	}

	return sugaredLogger.With(
		zap.String("logId", GetLogId(ctx)),
		zap.String("requestId", GetRequestId(ctx)),
		zap.String("module", env.GetAppName()),
		zap.String("localIp", env.GetLocalIP()),
		zap.String("uri", ctx.GetString(ContextKeyUri)),
	)
}

func Debug(ctx *gin.Context, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Debug(args...)
}

func DebugF(ctx *gin.Context, format string, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Debugf(format, args...)
}

func DebugW(ctx *gin.Context, msg string, keysAndValues ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Debugw(msg, keysAndValues...)
}

func Info(ctx *gin.Context, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Info(args...)
}

func InfoF(ctx *gin.Context, format string, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Infof(format, args...)
}

func InfoW(ctx *gin.Context, msg string, keysAndValues ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Infow(msg, keysAndValues...)
}

func Warn(ctx *gin.Context, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Warn(args...)
}

func WarnF(ctx *gin.Context, format string, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Warnf(format, args...)
}

func Error(ctx *gin.Context, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Error(args...)
}

func ErrorF(ctx *gin.Context, format string, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Errorf(format, args...)
}

func Panic(ctx *gin.Context, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Panic(args...)
}

func PanicF(ctx *gin.Context, format string, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Panicf(format, args...)
}

func Fatal(ctx *gin.Context, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Fatal(args...)
}

func Fatalf(ctx *gin.Context, format string, args ...interface{}) {
	if NoLog(ctx) {
		return
	}
	sLogger(ctx).Fatalf(format, args...)
}
