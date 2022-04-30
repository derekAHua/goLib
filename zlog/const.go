package zlog

// @Author: Derek
// @Description: Log constant.
// @Date: 2022/4/30 15:36
// @Version 1.0

// The prefix of log Name.
const (
	LogNameServer LogName = "server" // server 业务日志名字
	LogNameAccess         = "access" // access 日志文件名字
	LogNameModule         = "module" // module 日志文件名
	LogNameMysql          = "mysql"  // mysql 日志文件名
	LogNameRedis          = "redis"  // redis 日志文件名
	LogNameLua            = "lua"    // lua 日志文件名
	LogNameRMQ            = "rmq"    // rmq 日志文件名
	LogNameRpc            = "rpc"    // rpc 日志文件名
)
