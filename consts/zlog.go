package consts

// @Author: Derek
// @Description: Zlog constant identify.
// @Date: 2022/4/30 14:41
// @Version 1.0

const (
	TopicType = "_topic" // topic logo
)

// Log conf level.
const (
	LogLevelDefault = "INFO"
	LogLevelDebug   = "DEBUG"
	LogLevelInfo    = "INFO"
	LogLevelWarn    = "WARN"
	LogLevelError   = "ERROR"
	LogLevelFatal   = "FATAL"
)

// Log type
const (
	LogSuffixNormal    = ".log"
	LogSuffixWarnFatal = ".log.wf"
	LogStdout          = "stdout"
)

// Log util key
const (
	ContextKeyRequestId = "_requestId"
	ContextKeyLogId     = "_logId"
	ContextKeyNoLog     = "_noLogId"
	ContextKeyUri       = "_uri"
)

// Log header key
const (
	TraceHeaderKey = "requestId"
	LogIdHeaderKey = "logId"
)
