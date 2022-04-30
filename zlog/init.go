package zlog

import (
	"github.com/derekAHua/goLib/env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"sync"
)

type (
	LogLevel string

	LogConfig struct {
		Level  LogLevel `yaml:"level"`
		Stdout bool     `yaml:"stdout"`
	}

	loggerConfig struct {
		ZapLevel zapcore.Level

		Stdout bool
		Path   string
	}

	LogName string
)

// Only change by InitLog function.
var (
	logConfig = loggerConfig{
		ZapLevel: zapcore.InfoLevel,
		Stdout:   false,
		Path:     "./log",
	}

	onceLogConfig sync.Once
)

func InitLog(conf LogConfig) {
	onceLogConfig.Do(func() {
		logConfig.ZapLevel = getLogLevel(conf.Level)
		logConfig.Stdout = conf.Stdout
		logConfig.Path = env.GetLogDirPath()

		zapLogs := []LogName{LogNameServer, LogNameAccess, LogNameModule, LogNameMysql,
			LogNameRedis, LogNameLua, LogNameRMQ, LogNameRpc}

		for _, v := range zapLogs {
			if _, ok := mapZapLogger[v]; !ok {
				mapZapLogger[v] = newLogger(v).WithOptions(zap.AddCallerSkip(1))
			}
		}

		sugaredLogger = mapZapLogger[LogNameModule].Sugar()
	})
}
