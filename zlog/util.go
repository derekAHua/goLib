package zlog

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

// GetLogId return log ID.
// If logId is nil, will do GenId to get a new log ID and set it in context.
func GetLogId(ctx *gin.Context) string {
	if ctx == nil {
		return GenId()
	}

	// 从ctx中获取
	if logID := ctx.GetString(ContextKeyLogId); len(logID) > 0 {
		return logID
	}

	// 从header中获取
	var logID string
	if ctx.Request != nil && ctx.Request.Header != nil {
		logID = ctx.GetHeader(LogIdHeaderKey)
	}

	// 无logId，生成新的logId
	if logID == "" {
		logID = GenId()
	}

	ctx.Set(ContextKeyLogId, logID)
	return logID
}

// GetRequestId return request ID.
// If requestId is nil, will do GenId to get a new request ID and set it in context.
func GetRequestId(ctx *gin.Context) string {
	if ctx == nil {
		return GenId()
	}

	// 从ctx中获取
	if r := ctx.GetString(ContextKeyRequestId); r != "" {
		return r
	}

	// 从header中获取
	var requestId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		requestId = ctx.Request.Header.Get(TraceHeaderKey)
	}

	// 新生成
	if requestId == "" {
		requestId = GenId()
	}

	ctx.Set(ContextKeyRequestId, requestId)
	return requestId
}

func GenId() (requestId string) {
	u := uint64(time.Now().UnixNano())
	requestId = strconv.FormatUint(u&0x7FFFFFFF|0x80000000, 10)
	return requestId
}

func SetNoLogFlag(ctx *gin.Context) {
	ctx.Set(ContextKeyNoLog, true)
}

func NoLog(ctx *gin.Context) bool {
	if ctx == nil {
		return false
	}
	flag, ok := ctx.Get(ContextKeyNoLog)
	if ok && flag == true {
		return true
	}
	return false
}

func WithTopicField(logName string) zap.Field {
	return zap.String(TopicType, logName)
}
