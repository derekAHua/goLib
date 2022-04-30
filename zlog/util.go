package zlog

import (
	"github.com/derekAHua/goLib/consts"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

func GetLogId(ctx *gin.Context) string {
	if ctx == nil {
		return genRequestId()
	}

	// 从ctx中获取
	if logID := ctx.GetString(consts.ContextKeyLogId); logID != "" {
		return logID
	}

	// 从header中获取
	var logID string
	if ctx.Request != nil && ctx.Request.Header != nil {
		logID = ctx.GetHeader(consts.LogIdHeaderKey)
	}

	// 无logId，生成新的logId
	if logID == "" {
		logID = genRequestId()
	}

	ctx.Set(consts.ContextKeyLogId, logID)
	return logID
}

// GetRequestId 获取RequestId
func GetRequestId(ctx *gin.Context) string {
	if ctx == nil {
		return genRequestId()
	}

	// 从ctx中获取
	if r := ctx.GetString(consts.ContextKeyRequestId); r != "" {
		return r
	}

	// 从header中获取
	var requestId string
	if ctx.Request != nil && ctx.Request.Header != nil {
		requestId = ctx.Request.Header.Get(consts.TraceHeaderKey)
	}

	// 新生成
	if requestId == "" {
		requestId = genRequestId()
	}

	ctx.Set(consts.ContextKeyRequestId, requestId)
	return requestId
}

// Todo 获取分布式ID
func genRequestId() (requestId string) {
	u := uint64(time.Now().UnixNano())
	requestId = strconv.FormatUint(u&0x7FFFFFFF|0x80000000, 10)
	return requestId
}

func SetNoLogFlag(ctx *gin.Context) {
	ctx.Set(consts.ContextKeyNoLog, true)
}

func NoLog(ctx *gin.Context) bool {
	if ctx == nil {
		return false
	}
	flag, ok := ctx.Get(consts.ContextKeyNoLog)
	if ok && flag == true {
		return true
	}
	return false
}

func WithTopicField(logName LogName) zap.Field {
	return zap.String(consts.TopicType, string(logName))
}
