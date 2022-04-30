package zlog

import (
	"github.com/gin-gonic/gin"
	"strconv"
	"time"
)

// util key
const (
	ContextKeyRequestID = "requestId"
	ContextKeyLogID     = "logID"
	ContextKeyNoLog     = "_no_log"
	ContextKeyUri       = "_uri"
)

// header key
const (
	TraceHeaderKey      = "Uber-Trace-Id"
	LogIDHeaderKey      = "X_LogId"
	LogIDHeaderKeyLower = "x_logId"
)

func GetLogID(ctx *gin.Context) string {
	if ctx == nil {
		return genRequestId()
	}

	// 上次获取到的
	if logID := ctx.GetString(ContextKeyLogID); logID != "" {
		return logID
	}

	// 尝试从header中获取
	var logID string
	if ctx.Request != nil && ctx.Request.Header != nil {
		logID = ctx.GetHeader(LogIDHeaderKey)
		if logID == "" {
			logID = ctx.GetHeader(LogIDHeaderKeyLower)
		}
	}

	// 无法从上游获得，不展示logId，弱化logId
	if logID == "" {
		logID = genRequestId()
	}

	ctx.Set(ContextKeyLogID, logID)
	return logID
}

// GetRequestID 获取RequestId
func GetRequestID(ctx *gin.Context) string {
	if ctx == nil {
		return genRequestId()
	}

	// 从ctx中获取
	if r := ctx.GetString(ContextKeyRequestID); r != "" {
		return r
	}

	// 优先从header中获取
	var requestId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		requestId = ctx.Request.Header.Get(TraceHeaderKey)
	}

	// 新生成
	if requestId == "" {
		requestId = genRequestId()
	}

	ctx.Set(ContextKeyRequestID, requestId)
	return requestId
}

func genRequestId() (requestId string) {
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
