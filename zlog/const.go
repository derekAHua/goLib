package zlog

// @Author: Derek
// @Description: Log constant.
// @Date: 2022/4/30 15:36
// @Version 1.0

type LogName string

func (l LogName) ToString() string {
	return string(l)
}

// The prefix of log Name.
const (
	LogNameServer LogName = "server" // server 业务日志名字
	LogNameAccess         = "access" // access 日志文件名字
	LogNameMysql          = "mysql"  // mysql 日志文件名
	LogNameRedis          = "redis"  // redis 日志文件名
	LogNameLua            = "lua"    // lua 日志文件名
	LogNameRMQ            = "rmq"    // rmq 日志文件名
	LogNameRpc            = "rpc"    // rpc 日志文件名
	LogNameES             = "es"     // ES 日志文件名
)

const (
	TopicType = "_topic" // topic logo
)

// Log type.
const (
	LogSuffixNormal    = ".log"
	LogSuffixWarnFatal = ".log.wf"
	LogStdout          = "stdout"
)

// Log context key.
const (
	ContextKeyRequestId = "_requestId"
	ContextKeyLogId     = "_logId"
	ContextKeyNoLog     = "_noLogId"
	ContextKeyUri       = "_uri"
)

// Log header key.
const (
	TraceHeaderKey = "requestId"
	LogIdHeaderKey = "logId"
)
