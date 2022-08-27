package zlog

import (
	"fmt"
	"github.com/derekAHua/goLib/utils"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Field = zap.Field

// NewLogger return a new zapLogger.
func newLogger(name LogName) *zap.Logger {
	var (
		infoLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= logConfig.ZapLevel && lvl <= zapcore.InfoLevel
		})

		errorLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= logConfig.ZapLevel && lvl >= zapcore.WarnLevel
		})

		stdLevel = zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= logConfig.ZapLevel && lvl >= zapcore.DebugLevel
		})
	)

	if name == "" {
		panic("LogFile name must be set!")
	}

	var zapCore []zapcore.Core
	if logConfig.Stdout {
		c := zapcore.NewCore(getEncoder(), zapcore.AddSync(getLogWriter(name, LogStdout)), stdLevel)
		zapCore = append(zapCore, c)
	}

	zapCore = append(zapCore,
		zapcore.NewCore(getEncoder(), zapcore.AddSync(getLogWriter(name, LogSuffixNormal)), infoLevel))

	zapCore = append(zapCore,
		zapcore.NewCore(getEncoder(), zapcore.AddSync(getLogWriter(name, LogSuffixWarnFatal)), errorLevel),
	)

	core := zapcore.NewTee(zapCore...)
	filed := zap.Fields()
	caller := zap.AddCaller()
	development := zap.Development()
	return zap.New(core, filed, caller, development)
}

func getEncoder() zapcore.Encoder {
	timeEncoder := func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("2006-01-02 15:04:05.999999"))
	}

	encoderCfg := zapcore.EncoderConfig{
		LevelKey:       "level",
		TimeKey:        "time",
		CallerKey:      "file",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeCaller:   zapcore.ShortCallerEncoder, // 短路径编码器
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     timeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
	}
	return NewJsonEncoder(encoderCfg)
}

func getLogWriter(name LogName, loggerType string) zapcore.WriteSyncer {
	if loggerType == LogStdout {
		return zapcore.AddSync(os.Stdout)
	}

	logDir := logConfig.Path
	err := utils.MakeDirIfNo(logDir)
	if err != nil {
		panic(fmt.Errorf("create log dir '%s' error: %s", logDir, err))
	}

	filename := filepath.Join(strings.TrimSuffix(logDir, "/"), name.ToString()+loggerType)
	fd, err := os.OpenFile(filename, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
	if err != nil {
		panic("Open log file error: " + err.Error())
	}
	return zapcore.AddSync(fd)
}

func CloseLogger() {
	for _, logger := range mapZapLogger {
		if logger != nil {
			_ = logger.Sync()
		}
	}
}
